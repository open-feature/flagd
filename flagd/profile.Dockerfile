# Dockerfile with pprof profiler
# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS builder

WORKDIR /workspace
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG COMMIT
ARG DATE

# Copy source code
COPY flagd/ flagd
COPY core/ core

# Setup go workspace 
RUN go work init
RUN go work use ./flagd
RUN go work use ./core

# Go get dependencies
RUN cd flagd && go mod download

# Build with profiler
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" -o flagd-build flagd/main.go flagd/profiler.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/flagd-build .
USER 65532:65532

ENTRYPOINT ["/flagd-build"]
