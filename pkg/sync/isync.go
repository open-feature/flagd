package sync

import "context"

type ProviderArgs map[string]string

type Type int

// Type of the sync operation
const (
	// ALL - All flags of sync provider. This is the default if unset due to primitive default
	ALL Type = iota
	// ADD - Additional flags from sync provider
	ADD
	// UPDATE - Update for flag(s) previously provided
	UPDATE
	// DELETE - Delete for flag(s) previously provided
	DELETE
	RESYNC
)

func (t Type) String() string {
	switch t {
	case ALL:
		return "ALL"
	case ADD:
		return "ADD"
	case UPDATE:
		return "UPDATE"
	case DELETE:
		return "DELETE"
	case RESYNC:
		return "RESYNC"
	default:
		return "UNKNOWN"
	}
}

/*
ISync implementations watch for changes in the flag sources (HTTP backend, local file, K8s CRDs ...),fetch the latest
value and communicate to the Runtime with DataSync channel
*/
type ISync interface {
	// Sync is the contract between Runtime and sync implementation.
	// Note that, it is expected to return the first data sync as soon as possible to fill the store.
	Sync(ctx context.Context, dataSync chan<- DataSync) error

	ReSync(ctx context.Context, dataSync chan<- DataSync) error
}

// DataSync is the data contract between Runtime and sync implementations
type DataSync struct {
	FlagData string
	Source   string
	Type
}
