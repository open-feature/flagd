module github.com/open-feature/flagd/flagd-proxy/tests/loadtest

go 1.24

toolchain go1.24.1

require (
	buf.build/gen/go/open-feature-forking/flagd/grpc/go v1.6.0-20260204145352-a75813013224.1
	buf.build/gen/go/open-feature-forking/flagd/protocolbuffers/go v1.36.10-20260204145352-a75813013224.1
	google.golang.org/grpc v1.73.0
)

require (
	go.opentelemetry.io/otel v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.37.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
