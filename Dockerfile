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
# install the oapi-codegen binary
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.11.0

# Copy the go source
COPY main.go main.go
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY schemas/ schemas/

# Copy the code generation configs
COPY config/open_api_gen_config.yml open_api_gen_config.yml
COPY schemas/openapi/provider.yml provider.yml
COPY schemas/json-schema/flagd-definitions.json pkg/eval/flagd-definitions.json
# Generate OpenApi artifacts
RUN ${GOPATH}/bin/oapi-codegen --config=./open_api_gen_config.yml ./provider.yml
# Build
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -o flagd main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/flagd .
USER 65532:65532

ENTRYPOINT ["/flagd"]
