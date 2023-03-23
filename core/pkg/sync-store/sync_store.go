package store

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/runtime"
	isync "github.com/open-feature/flagd/core/pkg/sync"
	"go.uber.org/zap"
)

var (
	regCrd  *regexp.Regexp
	regFile *regexp.Regexp
)

func init() {
	regCrd = regexp.MustCompile("^core.openfeature.dev/")
	regFile = regexp.MustCompile("^file:")
}

type SyncStore struct {
	ctx          context.Context
	syncHandlers map[string]*syncHandler
	logger       *logger.Logger
	mu           *sync.RWMutex
	syncBuilder  SyncBuilderInterface
}

type syncHandler struct {
	subs       map[interface{}]storedChannels
	dataSync   chan isync.DataSync
	cancelFunc context.CancelFunc
	syncRef    isync.ISync
	mu         *sync.RWMutex
}

type storedChannels struct {
	errChan  chan error
	dataSync chan isync.DataSync
}

// NewSyncStore returns a new sync store
func NewSyncStore(ctx context.Context, logger *logger.Logger) *SyncStore {
	ss := SyncStore{
		ctx:          ctx,
		syncHandlers: map[string]*syncHandler{},
		logger:       logger,
		mu:           &sync.RWMutex{},
		syncBuilder:  &SyncBuilder{},
	}
	go ss.cleanup()
	return &ss
}

// FetchAllFlags returns a DataSync containing the full set of flag configurations from the SyncStore.
// This will either occur via triggering a resync, or through setting up a new subscription to the resource
func (s *SyncStore) FetchAllFlags(ctx context.Context, key interface{}, target string) (isync.DataSync, error) {
	s.logger.Debug(fmt.Sprintf("fetching all flags for target %s", target))
	dataSyncChan := make(chan isync.DataSync, 1)
	errChan := make(chan error, 1)
	s.mu.RLock()
	syncHandler, ok := s.syncHandlers[target]
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
func (s *SyncStore) RegisterSubscription(
	ctx context.Context,
	target string,
	key interface{},
	dataSync chan isync.DataSync,
	errChan chan error,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// is there a currently active subscription for this target?
	sh, ok := s.syncHandlers[target]
	if !ok {
		// we need to start a sync for this
		s.logger.Debug(
			fmt.Sprintf(
				"sync handler does not exist for target %s, registering syncHandler with sub %p",
				target,
				key,
			))
		sh = &syncHandler{
			dataSync: make(chan isync.DataSync),
			subs: map[interface{}]storedChannels{
				key: {
					errChan:  errChan,
					dataSync: dataSync,
				},
			},
			mu: &sync.RWMutex{},
		}
		s.syncHandlers[target] = sh
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
				s.logger.Debug(fmt.Sprintf("sync handler exists for target %s, triggering a resync", target))
				if err := sh.syncRef.ReSync(ctx, dataSync); err != nil {
					errChan <- err
				}
			}()
		}
	}
	// defer until context close to remove the key
	go func() {
		<-ctx.Done()
		s.mu.Lock()
		if s.syncHandlers[target] != nil && s.syncHandlers[target].subs != nil {
			s.logger.Debug(fmt.Sprintf("removing sync subscription due to context cancellation %p", key))
			delete(s.syncHandlers[target].subs, key)
		}
		s.mu.Unlock()
	}()
}

func (s *SyncStore) watchResource(target string) {
	s.logger.Debug(fmt.Sprintf("watching resource %s", target))
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()
	sh, ok := s.syncHandlers[target]
	if !ok {
		s.logger.Error(fmt.Sprintf("no sync handler exists for target %s", target))
		return
	}
	// this cancel is accessed by the cleanup method shutdown the listener + delete the syncHandler
	sh.cancelFunc = cancel
	go func() {
		<-ctx.Done()
		s.mu.Lock()
		delete(s.syncHandlers, target)
		s.mu.Unlock()
	}()
	// broadcast any data passed through the core channel to all subscribing channels
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d := <-sh.dataSync:
				sh.writeData(s.logger, d)
			}
		}
	}()
	// setup sync, if this fails an error is broadcasted, and the defer results in cleanup
	sync, err := s.syncBuilder.SyncFromURI(target, s.logger)
	if err != nil {
		s.logger.Error(fmt.Sprintf("unable to build sync from URI for target %s: %s", target, err.Error()))
		sh.writeError(s.logger, err)
		return
	}
	// init sync, if this fails an error is broadcasted, and the defer results in cleanup
	err = sync.Init(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("unable to initiate sync for target %s: %s", target, err.Error()))
		sh.writeError(s.logger, err)
		return
	}
	// sync ref is used to trigger a resync on a single channel when a new subscription is started
	// but the associated SyncHandler already exists, i.e. this function is not run
	sh.syncRef = sync
	err = sync.Sync(ctx, sh.dataSync)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error from sync for target %s: %s", target, err.Error()))
		sh.writeError(s.logger, err)
	}
}

func (h *syncHandler) writeError(logger *logger.Logger, err error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for k, ec := range h.subs {
		select {
		case ec.errChan <- err:
			continue
		default:
			logger.Error(fmt.Sprintf("unable to write error to channel for key %p", k))
		}
	}
}

func (h *syncHandler) writeData(logger *logger.Logger, data isync.DataSync) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for k, ds := range h.subs {
		select {
		case ds.dataSync <- data:
			continue
		default:
			logger.Error(fmt.Sprintf("unable to write data to channel for key %p", k))
		}
	}
}

func (s *SyncStore) cleanup() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(5 * time.Second):
			s.mu.Lock()
			for k, v := range s.syncHandlers {
				// delete any syncHandlers with 0 active subscriptions through cancelling its context
				s.logger.Debug(fmt.Sprintf("syncHandler for target %s has %d subscriptions", k, len(v.subs)))
				if len(v.subs) == 0 {
					s.logger.Debug(fmt.Sprintf("shutting down syncHandler %s", k))
					s.syncHandlers[k].cancelFunc()
				}
			}
			s.mu.Unlock()
		}
	}
}

type SyncBuilderInterface interface {
	SyncFromURI(uri string, logger *logger.Logger) (isync.ISync, error)
}

type SyncBuilder struct{}

// SyncFromURI builds an ISync interface from the input uri string
func (sb *SyncBuilder) SyncFromURI(uri string, logger *logger.Logger) (isync.ISync, error) {
	switch uriB := []byte(uri); {
	// filepath may be used for debugging, not recommended in deployment
	case regFile.Match(uriB):
		return runtime.NewFile(isync.SourceConfig{
			URI: regFile.ReplaceAllString(uri, ""),
		}, logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "filepath"),
			zap.String("target", "target"),
		)), nil
	case regCrd.Match(uriB):
		return runtime.NewK8s(uri, logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "kubernetes"),
		))
	}
	return nil, fmt.Errorf("unrecognized URI: %s", uri)
}
