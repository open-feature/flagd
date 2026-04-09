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

	// When true, the store scopes deletion to only the flagSetIds present in
	// this payload rather than wiping all flags for the source. This must be
	// explicitly opted-in per source via SourceConfig.IncrementalUpdates.
	// EXPERIMENTAL: this option may change or be removed in a future release.
	IncrementalUpdates bool
}

// SourceConfig is configuration option for flagd. This maps to startup parameter sources
type SourceConfig struct {
	URI      string `json:"uri"`
	Provider string `json:"provider"`

	AuthHeader string `json:"authHeader,omitempty"`
	CertPath   string `json:"certPath,omitempty"`
	TLS        bool   `json:"tls,omitempty"`
	ProviderID string `json:"providerID,omitempty"`
	Selector   string `json:"selector,omitempty"`
	Interval   uint32 `json:"interval,omitempty"`
	MaxMsgSize int    `json:"maxMsgSize,omitempty"`
	TimeoutS   int    `json:"timeoutS,omitempty"`

	// IncrementalUpdates opts this source into per-flagSetId scoped deletion.
	// When false (default), each update replaces all flags for the source.
	// When true, only flags matching the flagSetIds in the payload are replaced,
	// allowing flags from other flagSetIds to accumulate across updates.
	// EXPERIMENTAL: this option may change or be removed in a future release.
	// Note: flags from removed or renamed flagSetIds will not be automatically
	// cleaned up; a restart or explicit empty update for the old flagSetId is
	// required to purge them.
	IncrementalUpdates bool `json:"incrementalUpdates,omitempty"`

	OAuth *OAuthCredentialHandler `json:"oauth,omitempty"`
}

// OAuthCredentialHandler is a helper to manager OAuth 2.0 tokens, including re-loading of tokens.
type OAuthCredentialHandler struct {
	ClientID     string `json:"clientID,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	TokenURL     string `json:"tokenUrl,omitempty"`
	Folder       string `json:"folder,omitempty"`
	ReloadDelayS int    `json:"reloadDelayS,omitempty"`
}
