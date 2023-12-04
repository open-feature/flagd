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

// IManager defines the interface for the sync store
type IManager interface {
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

// Manager coordinates subscriptions by aggregating subscribers for the same target, and keeping them up to date
// for any updates that have happened for those targets.
type Manager struct {
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
func NewManager(ctx context.Context, logger *logger.Logger) *Manager {
	mgr := Manager{
		ctx:          ctx,
		multiplexers: map[string]*multiplexer{},
		logger:       logger,
		mu:           &sync.RWMutex{},
		syncBuilder:  &syncbuilder.SyncBuilder{},
	}
	go mgr.cleanup()
	return &mgr
}

// FetchAllFlags returns a DataSync containing the full set of flag configurations from the Manager.
// This will either occur via triggering a resync, or through setting up a new subscription to the resource
func (s *Manager) FetchAllFlags(ctx context.Context, key interface{}, target string) (isync.DataSync, error) {
	s.logger.Debug(fmt.Sprintf("fetching all flags for target %s", target))
	dataSyncChan := make(chan isync.DataSync, 1)
	errChan := make(chan error, 1)
	s.mu.RLock()
	syncHandler, ok := s.multiplexers[target]
	s.mu.RUnlock()
	if !ok {
		s.logger.Debug(fmt.Sprintf("sync handler does not exist for target %s, registering a new subscription", target))
		s.RegisterSubscription(ctx, target, key, dataSyncChan, errChan)
	} else {
		if syncHandler.syncRef == nil {
			return isync.DataSync{}, errors.New("sync ref not set")
		}
		go func() {
			s.logger.Debug(fmt.Sprintf("sync handler exists for target %s, triggering a resync", target))
			if err := syncHandler.syncRef.ReSync(ctx, dataSyncChan); err != nil {
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
func (s *Manager) RegisterSubscription(
	ctx context.Context,
	target string,
	key interface{},
	dataSync chan isync.DataSync,
	errChan chan error,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// is there a currently active subscription for this target?
	sh, ok := s.multiplexers[target]
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
			mu: &sync.RWMutex{},
		}
		go s.watchResource(target)
	} else {
		// register our sub in the map
		s.logger.Debug(fmt.Sprintf("registering sync subscription %p", key))
		sh.subs[key] = storedChannels{
			errChan:  errChan,
			dataSync: dataSync,
		}
		// access pointer + trigger resync passing the dataSync
		if sh.syncRef != nil {
			go func() {
				s.mu.RLock()
				defer s.mu.RUnlock()
				if _, ok := s.multiplexers[target]; ok {
					s.logger.Debug(fmt.Sprintf("sync handler exists for target %s, triggering a resync", target))
					if err := sh.syncRef.ReSync(ctx, dataSync); err != nil {
						errChan <- err
					}
				}
			}()
		}
	}
	// defer until context close to remove the key
	go func() {
		<-ctx.Done()
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.multiplexers[target] != nil && s.multiplexers[target].subs != nil {
			s.logger.Debug(fmt.Sprintf("removing sync subscription due to context cancellation %p", key))
			delete(s.multiplexers[target].subs, key)
		}
	}()
}

func (s *Manager) watchResource(target string) {
	s.logger.Debug(fmt.Sprintf("watching resource %s", target))
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()
	sh, ok := s.multiplexers[target]
	if !ok {
		s.logger.Error(fmt.Sprintf("no sync handler exists for target %s", target))
		return
	}
	// this cancel is accessed by the cleanup method shutdown the listener + delete the multiplexer
	sh.cancelFunc = cancel
	go func() {
		<-ctx.Done()
		s.mu.Lock()
		delete(s.multiplexers, target)
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
		sh.broadcastError(s.logger, err)
		return
	}
	// init sync, if this fails an error is broadcasted, and the defer results in cleanup
	err = syncSource.Init(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("unable to initiate sync for target %s: %s", target, err.Error()))
		sh.broadcastError(s.logger, err)
		return
	}
	// syncSource ref is used to trigger a resync on a single channel when a new subscription is started
	// but the associated SyncHandler already exists, i.e. this function is not run
	sh.syncRef = syncSource
	err = syncSource.Sync(ctx, sh.dataSync)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error from sync for target %s: %s", target, err.Error()))
		sh.broadcastError(s.logger, err)
	}
}

func (s *Manager) cleanup() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(5 * time.Second):
			s.mu.Lock()
			for k, v := range s.multiplexers {
				// delete any multiplexers with 0 active subscriptions through cancelling its context
				s.logger.Debug(fmt.Sprintf("multiplexer for target %s has %d subscriptions", k, len(v.subs)))
				if len(v.subs) == 0 {
					s.logger.Debug(fmt.Sprintf("shutting down multiplexer %s", k))
					s.multiplexers[k].cancelFunc()
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Manager) GetActiveSubscriptionsInt64() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	syncs := 0
	for _, v := range s.multiplexers {
		syncs += len(v.subs)
	}

	return int64(syncs)
}
