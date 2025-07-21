package sync

import (
	"context"

	"google.golang.org/protobuf/types/known/structpb"
)

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
	FlagData    string
	SyncContext *structpb.Struct
	Source      string
	Selector    string
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
