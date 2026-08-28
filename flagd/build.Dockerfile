# Main Dockerfile for flagd builds
# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.26-alpine@sha256:28d89ee9cc0ff9fec75c82ca201e6bf7fdf9a679d4b7b24dfa04f2bb766bb468 AS builder
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
    go work init ./core ./flagd && go mod download

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage a bind mount to the current directory to avoid having to copy the
# source code into the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,source=./core,target=./core \
    --mount=type=bind,source=./flagd,target=./flagd \
    if [ "${FIPS}" = on ]; then export GOFIPS140=v1.0.0 TAGS=fips140; else TAGS=; fi; \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -tags "${TAGS}" -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" -o /bin/flagd-build flagd/main.go

# Both tags together prove the variant built as intended: the toolchain adds
# fips140v1.0 only when GOFIPS140 is set, and fips140 is flagd's own.
RUN [ "${FIPS}" != on ] || go version -m /bin/flagd-build | grep -q -- '-tags=fips140,fips140v1.0'

# # Use distroless as minimal base image to package the manager binary
# # Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot@sha256:1c2c046bc09ed40fad370b599a0b1ae7987f55b01e247cf27a7c27cd97e5bbc7
WORKDIR /
COPY --from=builder /bin/flagd-build .
USER 65532:65532

ENTRYPOINT ["/flagd-build"]
