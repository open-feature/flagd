package runtime

import (
	"context"
	"errors"
	"testing"

	evalmock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
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

func TestRuntimeReadinessRequiresSuccessfulStateForEverySource(t *testing.T) {
	const (
		firstSource  = "file:/first.flagd.json"
		secondSource = "file:/second.flagd.json"
	)

	firstPayload := coresync.DataSync{Source: firstSource}
	secondPayload := coresync.DataSync{Source: secondSource}
	evaluatorError := errors.New("invalid flag configuration")

	ctrl := gomock.NewController(t)
	evaluator := evalmock.NewMockIEvaluator(ctrl)
	evaluator.EXPECT().SetState(firstPayload).Return(evaluatorError)
	evaluator.EXPECT().SetState(firstPayload).Return(nil)
	evaluator.EXPECT().SetState(secondPayload).Return(nil)

	syncService := &recordingSyncService{}
	runtime := Runtime{
		Evaluator:   evaluator,
		Logger:      logger.NewLogger(nil, false),
		SyncService: syncService,
		Syncs: []coresync.ISync{
			readySync{ready: true},
			readySync{ready: true},
		},
		sourceReadiness: map[string]bool{
			firstSource:  false,
			secondSource: false,
		},
	}

	require.False(t, runtime.isReady())

	runtime.updateAndEmit(firstPayload)
	require.False(t, runtime.isReady())

	runtime.updateAndEmit(firstPayload)
	require.False(t, runtime.isReady())

	runtime.updateAndEmit(secondPayload)
	require.True(t, runtime.isReady())
	require.Equal(t, []string{firstSource, secondSource}, syncService.emitted)
}
