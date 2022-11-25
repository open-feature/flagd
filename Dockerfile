# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.18-alpine AS builder

WORKDIR /workspace
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG COMMIT
ARG DATE
ARG GOBUILD
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY profiler.go profiler.go
COPY cmd/ cmd/
COPY pkg/ pkg/

# Check and include profiler in the build
RUN if [ "$GOBUILD" = "profile" ]; then \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" -o flagd main.go profiler.go ; \
   else \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" -o flagd main.go;  \
    fi

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/flagd .
USER 65532:65532

ENTRYPOINT ["/flagd"]
