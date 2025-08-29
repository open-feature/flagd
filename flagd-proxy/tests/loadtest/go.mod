module github.com/open-feature/flagd/flagd-proxy/tests/loadtest

go 1.23.0

toolchain go1.24.1

require (
	buf.build/gen/go/open-feature/flagd/grpc/go v1.5.1-20250529171031-ebdc14163473.2
	buf.build/gen/go/open-feature/flagd/protocolbuffers/go v1.36.6-20250529171031-ebdc14163473.1
	google.golang.org/grpc v1.73.0
)

require (
	go.opentelemetry.io/otel v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.37.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
