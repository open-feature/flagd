package runtime

import (
	"context"
	"errors"
	"testing"
	"time"

	evalmock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	coresync "github.com/open-feature/flagd/core/pkg/sync"
	flagsync "github.com/open-feature/flagd/flagd/pkg/service/flag-sync"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type readySync struct {
	ready bool
}

func (s readySync) Init(context.Context) error {
	return nil
}

func (s readySync) Sync(context.Context, chan<- coresync.DataSync) error {
	return nil
}

func (s readySync) ReSync(context.Context, chan<- coresync.DataSync) error {
	return nil
}

func (s readySync) IsReady() bool {
	return s.ready
}

type recordingSyncService struct {
	emitted []string
}

var _ flagsync.ISyncService = (*recordingSyncService)(nil)

func (s *recordingSyncService) Start(context.Context) error {
	return nil
}

func (s *recordingSyncService) Emit(source string) {
	s.emitted = append(s.emitted, source)
}

type controlledSync struct {
	payloads <-chan coresync.DataSync
	started  chan<- struct{}
}

func (s controlledSync) Init(context.Context) error {
	return nil
}

func (s controlledSync) Sync(ctx context.Context, dataSync chan<- coresync.DataSync) error {
	select {
	case s.started <- struct{}{}:
	case <-ctx.Done():
		return nil
	}

	for {
		select {
		case payload := <-s.payloads:
			select {
			case dataSync <- payload:
			case <-ctx.Done():
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s controlledSync) ReSync(context.Context, chan<- coresync.DataSync) error {
	return nil
}

func (s controlledSync) IsReady() bool {
	return true
}

type blockingEvaluationService struct {
	stop <-chan struct{}
}

var _ service.IFlagEvaluationService = blockingEvaluationService{}

func (s blockingEvaluationService) Serve(ctx context.Context, _ service.Configuration) error {
	select {
	case <-s.stop:
		return errors.New("test complete")
	case <-ctx.Done():
		return nil
	}
}

func (blockingEvaluationService) Notify(service.Notification) {}

func (blockingEvaluationService) Shutdown() {}

type blockingOfrepService struct{}

func (blockingOfrepService) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func requireSignal(t *testing.T, signal <-chan struct{}) {
	t.Helper()

	select {
	case <-signal:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for runtime signal")
	}
}

func TestRuntimeReadinessRequiresSuccessfulStateFromEveryProvider(t *testing.T) {
	const source = "file:/shared.flagd.json"

	payload := coresync.DataSync{Source: source}
	evaluatorError := errors.New("invalid flag configuration")

	ctrl := gomock.NewController(t)
	evaluator := evalmock.NewMockIEvaluator(ctrl)
	evaluator.EXPECT().SetState(payload).Return(evaluatorError)
	evaluator.EXPECT().SetState(payload).Return(nil)
	evaluator.EXPECT().SetState(payload).Return(nil)

	syncService := &recordingSyncService{}
	runtime := Runtime{
		Evaluator:   evaluator,
		Logger:      logger.NewLogger(nil, false),
		SyncService: syncService,
		Syncs: []coresync.ISync{
			readySync{ready: true},
			readySync{ready: true},
		},
		sourceReadiness: []bool{false, false},
	}

	require.False(t, runtime.isReady())

	runtime.updateAndEmit(payload, 0)
	require.False(t, runtime.isReady())

	runtime.updateAndEmit(payload, 0)
	require.False(t, runtime.isReady())

	runtime.updateAndEmit(payload, 1)
	require.True(t, runtime.isReady())
	require.Equal(t, []string{source, source}, syncService.emitted)
}

func TestRuntimeStartTracksProvidersWithSameSourceSeparately(t *testing.T) {
	const source = "file:/shared.flagd.json"

	payload := coresync.DataSync{Source: source}
	started := make(chan struct{}, 2)
	firstProviderData := make(chan coresync.DataSync)
	secondProviderData := make(chan coresync.DataSync)
	updated := make(chan struct{}, 2)
	stop := make(chan struct{})

	ctrl := gomock.NewController(t)
	evaluator := evalmock.NewMockIEvaluator(ctrl)
	evaluator.EXPECT().SetState(payload).DoAndReturn(func(coresync.DataSync) error {
		updated <- struct{}{}
		return nil
	}).Times(2)

	runtime := Runtime{
		Evaluator:         evaluator,
		Logger:            logger.NewLogger(nil, false),
		SyncService:       &recordingSyncService{},
		OfrepService:      blockingOfrepService{},
		EvaluationService: blockingEvaluationService{stop: stop},
		Syncs: []coresync.ISync{
			controlledSync{payloads: firstProviderData, started: started},
			controlledSync{payloads: secondProviderData, started: started},
		},
		sourceReadiness: []bool{false, false},
	}

	startErr := make(chan error, 1)
	go func() {
		startErr <- runtime.Start()
	}()

	requireSignal(t, started)
	requireSignal(t, started)

	firstProviderData <- payload
	requireSignal(t, updated)
	require.False(t, runtime.isReady())

	secondProviderData <- payload
	requireSignal(t, updated)
	require.True(t, runtime.isReady())

	close(stop)
	require.Error(t, <-startErr)
}
