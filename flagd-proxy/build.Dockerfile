# Main Dockerfile for flagd builds
# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.27-alpine@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc AS builder
# The toolchain determines which FIPS snapshots exist, so pin it.
ENV GOTOOLCHAIN=local

WORKDIR /src

ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG COMMIT
ARG DATE
# Set to "on" for the FIPS variant, built against the CMVP-certified Go
# Cryptographic Module v1.0.0 (certificate #5247).
ARG FIPS=off

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
    if [ "${FIPS}" = on ]; then export GOFIPS140=v1.0.0 TAGS=fips140; else TAGS=; fi; \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -tags "${TAGS}" -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" -o /bin/flagd-proxy flagd-proxy/main.go

# Both tags together prove the variant built as intended: the toolchain adds
# fips140v1.0 only when GOFIPS140 is set, and fips140 is flagd's own.
RUN [ "${FIPS}" != on ] || go version -m /bin/flagd-proxy | grep -q -- '-tags=fips140,fips140v1.0'

# # Use distroless as minimal base image to package the manager binary
# # Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot@sha256:1c2c046bc09ed40fad370b599a0b1ae7987f55b01e247cf27a7c27cd97e5bbc7
WORKDIR /
COPY --from=builder /bin/flagd-proxy .
USER 65532:65532

ENTRYPOINT ["/flagd-proxy"]
