# Main Dockerfile for flagd builds
# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /src

ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG COMMIT
ARG DATE

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage bind mounts to go.sum and go.mod to avoid having to copy them into
# the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=./core/go.mod,target=./core/go.mod \
    --mount=type=bind,source=./core/go.sum,target=./core/go.sum \
    --mount=type=bind,source=./flagd/go.mod,target=./flagd/go.mod \
    --mount=type=bind,source=./flagd/go.sum,target=./flagd/go.sum \
    --mount=type=bind,source=./flagd-proxy/go.mod,target=./flagd-proxy/go.mod \
    --mount=type=bind,source=./flagd-proxy/go.sum,target=./flagd-proxy/go.sum \
    go work init ./core ./flagd ./flagd-proxy && go mod download

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage a bind mount to the current directory to avoid having to copy the
# source code into the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,source=./core,target=./core \
    --mount=type=bind,source=./flagd,target=./flagd \
    --mount=type=bind,source=./flagd-proxy,target=./flagd-proxy \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" -o /bin/flagd-proxy flagd-proxy/main.go

# # Use distroless as minimal base image to package the manager binary
# # Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /bin/flagd-proxy .
USER 65532:65532

ENTRYPOINT ["/flagd-proxy"]
