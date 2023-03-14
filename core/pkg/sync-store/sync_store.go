package sync_store

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	isync "github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/flagd/core/pkg/sync/file"
	"github.com/open-feature/flagd/core/pkg/sync/kubernetes"
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
	mu           *sync.Mutex
	SyncBuilder  SyncBuilderInterface
}

type syncHandler struct {
	subs       map[interface{}]storedChannels
	dataSync   chan isync.DataSync
	cancelFunc context.CancelFunc
	syncRef    isync.ISync
}

type storedChannels struct {
	errChan  chan error
	dataSync chan isync.DataSync
}

func NewSyncStore(ctx context.Context, logger *logger.Logger) SyncStore {
	return SyncStore{
		ctx:          ctx,
		syncHandlers: map[string]*syncHandler{},
		logger:       logger,
		mu:           &sync.Mutex{},
		SyncBuilder:  &SyncBuilder{},
	}
}

func (s *SyncStore) FetchAllFlags(ctx context.Context, key interface{}, target string) (isync.DataSync, error) {
	dataSyncChan := make(chan isync.DataSync, 1)
	errChan := make(chan error, 1)
	s.mu.Lock()
	_, ok := s.syncHandlers[target]
	if !ok {
		s.mu.Unlock()
		s.RegisterSubscription(ctx, target, key, dataSyncChan, errChan)
	} else {
		go s.syncHandlers[target].syncRef.ReSync(ctx, dataSyncChan)
		s.mu.Unlock()
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

func (s *SyncStore) RegisterSubscription(ctx context.Context, target string, key interface{}, dataSync chan isync.DataSync, errChan chan error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// is there a currently active subscription for this target?
	_, ok := s.syncHandlers[target]
	if !ok {
		// we need to start a sync for this
		syncHandler := syncHandler{
			dataSync: make(chan isync.DataSync),
			subs: map[interface{}]storedChannels{
				key: {
					errChan:  errChan,
					dataSync: dataSync,
				},
			},
		}
		s.syncHandlers[target] = &syncHandler
		fmt.Println("here", s.syncHandlers[target])
		go s.watchResource(s.ctx, target)
	} else {
		// register our sub in the map
		s.syncHandlers[target].subs[key] = storedChannels{
			errChan:  errChan,
			dataSync: dataSync,
		}

		// access pointer + trigger resync passing the dataSync
		if s.syncHandlers[target].syncRef != nil {
			go s.syncHandlers[target].syncRef.ReSync(ctx, dataSync)
		}

	}
	// defer until context close to remove the key
	go func() {
		<-ctx.Done()
		s.mu.Lock()
		if s.syncHandlers[target] != nil && s.syncHandlers[target].subs != nil {
			delete(s.syncHandlers[target].subs, key)
		}
		s.mu.Unlock()
	}()
	return nil
}

func (s *SyncStore) watchResource(ctx context.Context, target string) {
	ctx, cancel := context.WithCancel(ctx)
	s.syncHandlers[target].cancelFunc = cancel
	defer cancel()

	go func() {
		<-ctx.Done()
		fmt.Println("deleting")
		s.mu.Lock()
		delete(s.syncHandlers, target)
		s.mu.Unlock()
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d := <-s.syncHandlers[target].dataSync:
				s.mu.Lock()
				for _, ds := range s.syncHandlers[target].subs {
					ds.dataSync <- d
				}
				s.mu.Unlock()
			}
		}
	}()

	sync, err := s.SyncBuilder.SyncFromURI(target, *s.logger)
	if err != nil {
		s.mu.Lock()
		for _, ec := range s.syncHandlers[target].subs {
			ec.errChan <- err
		}
		s.mu.Unlock()
		return
	}
	err = sync.Init(ctx)
	if err != nil {
		s.mu.Lock()
		for _, ec := range s.syncHandlers[target].subs {
			ec.errChan <- err
		}
		s.mu.Unlock()
		return
	}

	s.syncHandlers[target].syncRef = sync

	err = sync.Sync(ctx, s.syncHandlers[target].dataSync)
	if err != nil {
		s.mu.Lock()
		for _, ec := range s.syncHandlers[target].subs {
			ec.errChan <- err
		}
		s.mu.Unlock()
	}
}

func (s *SyncStore) Cleanup() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(5 * time.Second):
			s.mu.Lock()
			for k, v := range s.syncHandlers {
				if len(v.subs) == 0 {
					s.syncHandlers[k].cancelFunc()
				}
			}
			s.mu.Unlock()
		}
	}
}

type SyncBuilderInterface interface {
	SyncFromURI(uri string, logger logger.Logger) (isync.ISync, error)
}

type SyncBuilder struct{}

func (sb *SyncBuilder) SyncFromURI(uri string, logger logger.Logger) (isync.ISync, error) {
	switch uriB := []byte(uri); {
	case regFile.Match(uriB):
		return &file.Sync{
			URI: regFile.ReplaceAllString(uri, ""),
			Logger: logger.WithFields(
				zap.String("component", "sync"),
				zap.String("sync", "filepath"),
				zap.String("target", "target"),
			),
			Mux: &sync.RWMutex{},
		}, nil
	case regCrd.Match(uriB):
		reader, dynamic, err := kubernetes.GetClients()
		if err != nil {
			return nil, err
		}
		return kubernetes.NewK8sSync(
			logger.WithFields(
				zap.String("component", "sync"),
				zap.String("sync", "kubernetes"),
			),
			regCrd.ReplaceAllString(uri, ""),
			reader,
			dynamic,
		), nil
	}
	return nil, fmt.Errorf("unrecognized URI: %s", uri)
}
