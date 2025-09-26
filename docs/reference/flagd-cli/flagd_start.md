<!-- markdownlint-disable-file -->
<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->
## flagd start

Start flagd

```
flagd start [flags]
```

### Options

```
  -H, --context-from-header stringToString   add key-value pairs to map header values to context values, where key is Header name, value is context key (default [])
  -X, --context-value stringToString         add arbitrary key value pairs to the flag evaluation context (default [])
  -C, --cors-origin strings                  CORS allowed origins, * will allow all origins
      --disable-sync-metadata                Disables the getMetadata endpoint of the sync service. Defaults to false, but will default to true in later versions.
  -h, --help                                 help for start
  -z, --log-format string                    Set the logging format, e.g. console or json (default "console")
  -m, --management-port int32                Port for management operations (default 8014)
  -t, --metrics-exporter string              Set the metrics exporter. Default(if unset) is Prometheus. Can be override to otel - OpenTelemetry metric exporter. Overriding to otel require otelCollectorURI to be present
  -r, --ofrep-port int32                     ofrep service port (default 8016)
  -A, --otel-ca-path string                  tls certificate authority path to use with OpenTelemetry collector
  -D, --otel-cert-path string                tls certificate path to use with OpenTelemetry collector
  -o, --otel-collector-uri string            Set the grpc URI of the OpenTelemetry collector for flagd runtime. If unset, the collector setup will be ignored and traces will not be exported.
  -K, --otel-key-path string                 tls key path to use with OpenTelemetry collector
  -I, --otel-reload-interval duration        how long between reloading the otel tls certificate from disk (default 1h0m0s)
  -p, --port int32                           Port to listen on (default 8013)
      --selector-fallback-key string         Fallback key for a selector expression that does not contain an '='. Defaults to 'source'.
  -c, --server-cert-path string              Server side tls certificate path
  -k, --server-key-path string               Server side tls key path
  -d, --socket-path string                   Flagd unix socket path. With grpc the evaluations service will become available on this address. With http(s) the grpc-gateway proxy will use this address internally.
  -s, --sources string                       JSON representation of an array of SourceConfig objects. This object contains 2 required fields, uri (string) and provider (string). Documentation for this object: https://flagd.dev/reference/sync-configuration/#source-configuration
      --stream-deadline duration             Set a server-side deadline for flagd sync and event streams (default 0, means no deadline).
  -g, --sync-port int32                      gRPC Sync port (default 8015)
  -e, --sync-socket-path string              Flagd sync service socket path. With grpc the sync service will be available on this address.
  -f, --uri .yaml/.yml/.json                 Set a sync provider uri to read data from, this can be a filepath, URL (HTTP and gRPC), FeatureFlag custom resource, or GCS or Azure Blob. When flag keys are duplicated across multiple providers the merge priority follows the index of the flag arguments, as such flags from the uri at index 0 take the lowest precedence, with duplicated keys being overwritten by those from the uri at index 1. Please note that if you are using filepath, flagd only supports files with .yaml/.yml/.json extension.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.agent.yaml)
  -x, --debug           verbose logging
```

### SEE ALSO

* [flagd](flagd.md)	 - Flagd is a simple command line tool for fetching and presenting feature flags to services. It is designed to conform to Open Feature schema for flag definitions.

