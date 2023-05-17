package store

import (
	"context"

	isync "github.com/open-feature/flagd/core/pkg/sync"
)

// ISyncStore defines the interface for the sync store
type ISyncStore interface {
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
