# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.18-alpine AS builder

WORKDIR /workspace
ARG TARGETOS
ARG TARGETARCH
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
# install buf and protoc binaries
RUN go install github.com/bufbuild/buf/cmd/buf@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Copy the go source
COPY main.go main.go
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY schemas/ schemas/

# Copy the code generation configs
COPY schemas/protobuf schemas/protobuf
COPY schemas/json/flagd-definitions.json pkg/eval/flagd-definitions.json
# Generate http/grpc stubs
RUN cd schemas/protobuf && ${GOPATH}/bin/buf generate --template buf.gen.go-server.yaml && cd ../..
# Build
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -o flagd main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/flagd .
USER 65532:65532

ENTRYPOINT ["/flagd"]
