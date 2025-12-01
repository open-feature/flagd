package builder

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/open-feature/flagd/core/pkg/sync"
)

// ParseSources parse a json formatted SourceConfig array string and performs validations on the content
func ParseSources(sourcesFlag string) ([]sync.SourceConfig, error) {
	syncProvidersParsed := []sync.SourceConfig{}

	if err := json.Unmarshal([]byte(sourcesFlag), &syncProvidersParsed); err != nil {
		return syncProvidersParsed, fmt.Errorf("error parsing sync providers: %w", err)
	}
	for _, sp := range syncProvidersParsed {
		if sp.URI == "" {
			return syncProvidersParsed, errors.New("sync provider argument parse: uri is a required field")
		}
		if sp.Provider == "" {
			return syncProvidersParsed, errors.New("sync provider argument parse: provider is a required field")
		}
	}
	return syncProvidersParsed, nil
}

// ParseSyncProviderURIs uri flag based sync sources to SourceConfig array. Replaces uri prefixes where necessary to
// derive SourceConfig
func ParseSyncProviderURIs(uris []string) ([]sync.SourceConfig, error) {
	syncProvidersParsed := []sync.SourceConfig{}

	for _, uri := range uris {
		switch uriB := []byte(uri); {
		case regFile.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.SourceConfig{
				URI:      regFile.ReplaceAllString(uri, ""),
				Provider: syncProviderFile,
			})
		case regCrd.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.SourceConfig{
				URI:      regCrd.ReplaceAllString(uri, ""),
				Provider: syncProviderKubernetes,
			})
		case regURL.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.SourceConfig{
				URI:      uri,
				Provider: syncProviderHTTP,
			})
		case regGRPC.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.SourceConfig{
				URI:      regGRPC.ReplaceAllString(uri, ""),
				Provider: syncProviderGrpc,
			})
		case regGRPCSecure.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.SourceConfig{
				URI:      regGRPCSecure.ReplaceAllString(uri, ""),
				Provider: syncProviderGrpc,
				TLS:      true,
			})
		case regGRPCCustomResolver.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.SourceConfig{
				URI:      uri,
				Provider: syncProviderGrpc,
			})
		case regGcs.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.SourceConfig{
				URI:      uri,
				Provider: syncProviderGcs,
			})
		case regAzblob.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.SourceConfig{
				URI:      uri,
				Provider: syncProviderAzblob,
			})
		case regS3.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.SourceConfig{
				URI:      uri,
				Provider: syncProviderS3,
			})
		default:
			return syncProvidersParsed, fmt.Errorf("invalid sync uri argument: %s, must start with 'file:', "+
				"'http(s)://', 'grpc(s)://', 'gs://', 'azblob://' or 'core.openfeature.dev'", uri)
		}
	}
	return syncProvidersParsed, nil
}
