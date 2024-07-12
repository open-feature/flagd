package sync

import (
	"context"
)

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
	default:
		return "UNKNOWN"
	}
}

/*
ISync implementations watch for changes in the flag sources (HTTP backend, local file, K8s CRDs ...),fetch the latest
value and communicate to the Runtime with DataSync channel
*/
type ISync interface {
	// Init is used by the sync provider to initialize its data structures and external dependencies.
	Init(ctx context.Context) error

	// Sync is the contract between Runtime and sync implementation.
	// Note that, it is expected to return the first data sync as soon as possible to fill the store.
	Sync(ctx context.Context, dataSync chan<- DataSync) error

	// ReSync is used to fetch the full flag configuration from the sync
	// This method should trigger an ALL sync operation then exit
	ReSync(ctx context.Context, dataSync chan<- DataSync) error

	// IsReady shall return true if the provider is ready to communicate with the Runtime
	IsReady() bool
}

// DataSync is the data contract between Runtime and sync implementations
type DataSync struct {
	FlagData string
	Source   string
	Selector string
	Type
}

// SourceConfig is configuration option for flagd. This maps to startup parameter sources
type SourceConfig struct {
	URI      string `json:"uri"`
	Provider string `json:"provider"`

	BearerToken string `json:"bearerToken,omitempty"`
	AuthHeader  string `json:"authHeader,omitempty"`
	CertPath    string `json:"certPath,omitempty"`
	TLS         bool   `json:"tls,omitempty"`
	ProviderID  string `json:"providerID,omitempty"`
	Selector    string `json:"selector,omitempty"`
	Interval    uint32 `json:"interval,omitempty"`
	MaxMsgSize  int    `json:"maxMsgSize,omitempty"`
}
