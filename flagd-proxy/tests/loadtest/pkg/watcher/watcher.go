package watcher

import (
	"context"
	"fmt"
	"io"
	"time"

	"buf.build/gen/go/open-feature-forking/flagd/grpc/go/sync/v1/syncv1grpc"
	syncv1Types "buf.build/gen/go/open-feature-forking/flagd/protocolbuffers/go/sync/v1"
)

const (
	timeoutSeconds = 10
)

type Watcher struct {
	client syncv1grpc.FlagSyncServiceClient
	//nolint:staticcheck
	Stream     chan syncv1Types.SyncState
	Ready      chan struct{}
	targetFile string
}

func NewWatcher(client syncv1grpc.FlagSyncServiceClient, target string) *Watcher {
	return &Watcher{
		//nolint:staticcheck
		Stream:     make(chan syncv1Types.SyncState, 1),
		client:     client,
		Ready:      make(chan struct{}),
		targetFile: target,
	}
}

//nolint:staticcheck
func (w *Watcher) StartWatcher(ctx context.Context) error {
	stream, err := w.client.SyncFlags(ctx, &syncv1Types.SyncFlagsRequest{
		Selector: fmt.Sprintf("file:%s", w.targetFile),
	})
	if err != nil {
		return fmt.Errorf("unable to create stream: %w", err)
	}

	ready := false
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("unable to read payload from stream: %w", err)
		}
		w.Stream <- msg.State
		if !ready {
			ready = true
			close(w.Ready)
		}
	}
}

func (w *Watcher) Wait() error {
	w.drainChan()
	select {
	case <-time.After(timeoutSeconds * time.Second):
		return fmt.Errorf("timeout out after %d", timeoutSeconds)
	case <-w.Stream:
		return nil
	}
}

func (w *Watcher) drainChan() {
	for {
		select {
		case <-w.Stream:
		default:
			return
		}
	}
}
