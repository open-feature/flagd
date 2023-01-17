package sync

import "context"

/*
ISync implementations watch for changes in the flag sources (HTTP backend, local file, K8s CRDs ...),fetch the latest
value and communicate to the Runtime with DataSync channel
*/

type ProviderArgs map[string]string

type ISync interface {
	// Sync is the contract between Runtime and sync implementation.
	// Note that, it is expected to return the first data sync as soon as possible to fill the store.
	Sync(ctx context.Context, dataSync chan<- DataSync) error
}

// DataSync is the data contract between Runtime and sync implementations
type DataSync struct {
	FlagData string
	Source   string
}
