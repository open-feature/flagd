//nolint:contextcheck
package subscriptions

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	isync "github.com/open-feature/flagd/core/pkg/sync"
	syncbuilder "github.com/open-feature/flagd/core/pkg/sync/builder"
)

// Manager defines the interface for the subscription management
type Manager interface {
	FetchAllFlags(
		ctx context.Context,
		key interface{},
		target string,
	) (isync.DataSync, error)
	RegisterSubscription(
		ctx context.Context,
		target string,
		key interface{},
		dataSync chan isync.DataSync,
		errChan chan error,
	)

	// metrics hooks
	GetActiveSubscriptionsInt64() int64
}

// Coordinator coordinates subscriptions by aggregating subscribers for the same target, and keeping them up to date
// for any updates that have happened for those targets.
type Coordinator struct {
	ctx          context.Context
	multiplexers map[string]*multiplexer
	logger       *logger.Logger
	mu           *sync.RWMutex
	syncBuilder  syncbuilder.ISyncBuilder
}

type storedChannels struct {
	errChan  chan error
	dataSync chan isync.DataSync
}

// NewManager returns a new subscription manager
func NewManager(ctx context.Context, logger *logger.Logger) *Coordinator {
	mgr := Coordinator{
		ctx:          ctx,
		multiplexers: map[string]*multiplexer{},
		logger:       logger,
		mu:           &sync.RWMutex{},
		syncBuilder:  syncbuilder.NewSyncBuilder(),
	}
	go mgr.cleanup()
	return &mgr
}

// FetchAllFlags returns a DataSync containing the full set of flag configurations from the Coordinator.
// This will either occur via triggering a resync, or through setting up a new subscription to the resource
func (s *Coordinator) FetchAllFlags(ctx context.Context, key interface{}, target string) (isync.DataSync, error) {
	s.logger.Debug(fmt.Sprintf("fetching all flags for target %s", target))
	dataSyncChan := make(chan isync.DataSync, 1)
	errChan := make(chan error, 1)
	s.mu.RLock()
	syncHandler, ok := s.multiplexers[target]
	// syncRef is written by watchResource under s.mu, so read it while we still hold the lock
	var syncRef isync.ISync
	if ok {
		syncRef = syncHandler.syncRef
	}
	s.mu.RUnlock()
	if !ok {
		s.logger.Debug(fmt.Sprintf("sync handler does not exist for target %s, registering a new subscription", target))
		s.RegisterSubscription(ctx, target, key, dataSyncChan, errChan)
	} else {
		if syncRef == nil {
			return isync.DataSync{}, errors.New("sync ref not set")
		}
		go func() {
			s.logger.Debug(fmt.Sprintf("sync handler exists for target %s, triggering a resync", target))
			if err := syncRef.ReSync(ctx, dataSyncChan); err != nil {
				errChan <- err
			}
		}()
	}

	select {
	case data := <-dataSyncChan:
		return data, nil
	case err := <-errChan:
		return isync.DataSync{}, err
	case <-time.After(5 * time.Second):
		return isync.DataSync{}, errors.New("fetching all flags timed out after 5 seconds")
	}
}

// RegisterSubscription starts a new subscription to the target resource.
// Once the subscription is set an ALL sync event will be received via the DataSync chan.
func (s *Coordinator) RegisterSubscription(
	ctx context.Context,
	target string,
	key interface{},
	dataSync chan isync.DataSync,
	errChan chan error,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// is there a currently active subscription for this target?
	// a multiplexer whose watcher has stopped can never deliver again, so treat it as absent (#2030)
	sh, ok := s.multiplexers[target]
	if ok && sh.isDead() {
		s.logger.Debug(fmt.Sprintf("multiplexer for target %s is no longer watched, replacing it", target))
		ok = false
	}
	if !ok {
		// we need to start a sync for this
		s.logger.Debug(
			fmt.Sprintf(
				"sync handler does not exist for target %s, registering multiplexer with sub %p",
				target,
				key,
			))
		s.multiplexers[target] = &multiplexer{
			dataSync: make(chan isync.DataSync),
			subs: map[interface{}]storedChannels{
				key: {
					errChan:  errChan,
					dataSync: dataSync,
				},
			},
		}
		go s.watchResource(target)
	} else {
		s.logger.Debug(fmt.Sprintf("registering sync subscription %p", key))
		sh.addSub(key, storedChannels{errChan: errChan, dataSync: dataSync})
		// the goroutine takes neither lock: the subscriber channel is unbuffered and can stall
		if syncRef := sh.syncRef; syncRef != nil {
			go func() {
				s.logger.Debug(fmt.Sprintf("sync handler exists for target %s, triggering a resync", target))
				if err := syncRef.ReSync(ctx, dataSync); err != nil {
					errChan <- err
				}
			}()
		}
	}
	// defer until context close to remove the key, from the multiplexer we actually joined:
	// a rebuild may have replaced the entry for this target by then (#2030)
	registered := s.multiplexers[target]
	go func() {
		<-ctx.Done()
		s.mu.Lock()
		defer s.mu.Unlock()
		s.logger.Debug(fmt.Sprintf("removing sync subscription due to context cancellation %p", key))
		registered.removeSub(key)
	}()
}

func (s *Coordinator) watchResource(target string) {
	s.logger.Debug(fmt.Sprintf("watching resource %s", target))
	ctx, cancel := context.WithCancel(s.ctx)
	s.mu.Lock()
	sh, ok := s.multiplexers[target]
	if !ok {
		s.mu.Unlock()
		cancel() // bare: no multiplexer exists for this target, so there is nothing to mark
		s.logger.Error(fmt.Sprintf("no sync handler exists for target %s", target))
		return
	}
	// keeps the Coordinator.mu -> multiplexer.mu order
	sh.watchedBy(cancel)
	s.mu.Unlock()
	var broadcastErr error
	defer func() {
		// one hold of Coordinator.mu for the whole teardown, so a RegisterSubscription cannot
		// land in between and attach to a multiplexer about to die (#2030)
		s.mu.Lock()
		sh.kill()
		// only our own entry; a later subscription may already have replaced it
		if s.multiplexers[target] == sh {
			delete(s.multiplexers, target)
		}
		// still under the lock: whoever registered before this is still a subscriber here
		if broadcastErr != nil {
			sh.broadcastError(s.logger, broadcastErr)
		}
		s.mu.Unlock()
	}()
	// broadcast any data passed through the core channel to all subscribing channels
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d := <-sh.dataSync:
				sh.broadcastData(s.logger, d)
			}
		}
	}()
	// setup sync, if this fails an error is broadcasted, and the defer results in cleanup
	syncSource, err := s.syncBuilder.SyncFromURI(target, s.logger)
	if err != nil {
		s.logger.Error(fmt.Sprintf("unable to build sync from URI for target %s: %s", target, err.Error()))
		broadcastErr = err
		return
	}
	// init sync, if this fails an error is broadcasted, and the defer results in cleanup
	err = syncSource.Init(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("unable to initiate sync for target %s: %s", target, err.Error()))
		broadcastErr = err
		return
	}
	// syncSource ref is used to trigger a resync on a single channel when a new subscription is started
	// but the associated SyncHandler already exists, i.e. this function is not run.
	// written under s.mu because the readers hold it
	s.mu.Lock()
	sh.syncRef = syncSource
	s.mu.Unlock()
	err = syncSource.Sync(ctx, sh.dataSync)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error from sync for target %s: %s", target, err.Error()))
		broadcastErr = err
	}
}

func (s *Coordinator) cleanup() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(5 * time.Second):
			s.mu.Lock()
			for k, v := range s.multiplexers {
				// reap any multiplexer with 0 active subscriptions; kill, never a bare cancel (#2030)
				s.logger.Debug(fmt.Sprintf("multiplexer for target %s has %d subscriptions", k, v.subCount()))
				if v.subCount() == 0 {
					s.logger.Debug(fmt.Sprintf("shutting down multiplexer %s", k))
					s.multiplexers[k].kill()
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Coordinator) GetActiveSubscriptionsInt64() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	syncs := 0
	for _, v := range s.multiplexers {
		syncs += v.subCount()
	}

	return int64(syncs)
}
