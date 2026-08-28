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
# Build variant. The defaults produce the standard build; the FIPS variant is
# built with GOFIPS140=v1.0.0 and GO_BUILD_TAGS=fips140, which selects the
# CMVP-certified Go Cryptographic Module (certificate #5247) and makes the
# binary require it at startup. Use the literal version, not the
# "certified"/"inprocess" aliases, which vary by toolchain.
ARG GOFIPS140=off
ARG GO_BUILD_TAGS=

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
    CGO_ENABLED=0 GOFIPS140=${GOFIPS140} GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -tags "${GO_BUILD_TAGS}" -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" -o /bin/flagd-build flagd/main.go

# Fail closed if a FIPS build was requested but the settings did not reach the binary.
RUN if [ "${GOFIPS140}" != "off" ]; then \
      go version -m /bin/flagd-build | grep -q 'GOFIPS140=v1\.0\.0' \
        || (echo "ERROR: /bin/flagd-build is not a FIPS 140-3 build" && exit 1); \
      go version -m /bin/flagd-build | grep -E 'GOFIPS140=|-tags='; \
    fi

# # Use distroless as minimal base image to package the manager binary
# # Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot@sha256:1c2c046bc09ed40fad370b599a0b1ae7987f55b01e247cf27a7c27cd97e5bbc7
WORKDIR /
COPY --from=builder /bin/flagd-build .
USER 65532:65532

ENTRYPOINT ["/flagd-build"]
