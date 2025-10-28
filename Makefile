PHONY: .docker-build .build .run .mockgen
PREFIX=/usr/local
PUBLIC_JSON_SCHEMA_DIR=docs/schema/v0/
ALL_GO_MOD_DIRS := $(shell find . -path ./test/integration -prune -o -type f -name 'go.mod' -exec dirname {} \; | sort)

FLAGD_DEV_NAMESPACE ?= flagd-dev
ZD_TEST_NAMESPACE_FLAGD_PROXY ?= flagd-proxy-zd-test
ZD_TEST_NAMESPACE ?= flagd-zd-test
ZD_CLIENT_IMG ?= zd-client:latest
FLAGD_PROXY_IMG ?= flagd-proxy:latest
FLAGD_PROXY_IMG_ZD ?= flagd-proxy:zd

DOCS_DIR ?= docs

workspace-init: workspace-clean
	go work init
	$(foreach module, $(ALL_GO_MOD_DIRS), go work use $(module);)

workspace-update:
	$(foreach module, $(ALL_GO_MOD_DIRS), go work use $(module);)

workspace-clean:
	rm -rf go.work

guard-%:
	@ if [ "${${*}}" = "" ]; then \
        echo "Environment variable $* not set"; \
        exit 1; \
    fi
docker-build: # default to flagd
	make docker-build-flagd
docker-push: # default to flagd
	make docker-push-flagd
docker-build-flagd:
	docker buildx build --build-arg=VERSION="$$(git describe --tags --abbrev=0)" --build-arg=COMMIT="$$(git rev-parse --short HEAD)" --build-arg DATE="$$(date +%FT%TZ)" --platform="linux/arm64" -t ${IMG} -f flagd/build.Dockerfile .
docker-push-flagd:
	docker buildx build --push --build-arg=VERSION="$$(git describe --tags --abbrev=0)" --build-arg=COMMIT="$$(git rev-parse --short HEAD)" --build-arg DATE="$$(date +%FT%TZ)" --platform="linux/ppc64le,linux/s390x,linux/amd64,linux/arm64" -t ${IMG} -f flagd/build.Dockerfile .
build: workspace-init # default to flagd
	make build-flagd
build-flagd:
	go build -ldflags "-X main.version=dev -X main.commit=$$(git rev-parse --short HEAD) -X main.date=$$(date +%FT%TZ)" -o ./bin/flagd ./flagd
.PHONY: test
test: test-core test-flagd test-flagd-proxy
test-core:
	go test -race -covermode=atomic -cover -short ./core/pkg/... -coverprofile=core-coverage.out
test-flagd:
	go test -race -covermode=atomic -cover -short ./flagd/pkg/... -coverprofile=flagd-coverage.out
test-flagd-proxy:
	go test -race -covermode=atomic -cover -short ./flagd-proxy/pkg/... -coverprofile=flagd-proxy-coverage.out
flagd-benchmark-test:
	go test -bench=Bench -short -benchtime=5s -benchmem ./core/... | tee benchmark.txt
flagd-integration-test-harness:
# target used to start a locally built flagd with the e2e flags
	cd flagd; go run main.go start -f file:../test-harness/flags/testing-flags.json -f file:../test-harness/flags/custom-ops.json -f file:../test-harness/flags/evaluator-refs.json -f file:../test-harness/flags/zero-flags.json -f file:../test-harness/flags/edge-case-flags.json
flagd-integration-test: # dependent on flagd-e2e-test-harness if not running in github actions
	go test -count=1 -cover ./test/integration $(ARGS)
run: # default to flagd
	make run-flagd
run-flagd:
	cd flagd; go run main.go start -f file:../config/samples/example_flags.flagd.json 
run-flagd-selector-demo:
	cd flagd; go run main.go start -f file:../config/samples/example_flags.flagd.json -f file:../config/samples/example_flags.flagd.2.json
install:
	cp systemd/flagd.service /etc/systemd/system/flagd.service
	mkdir -p /etc/flagd
	cp systemd/flags.json /etc/flagd/flags.json
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp bin/flagd $(DESTDIR)$(PREFIX)/bin/flagd
	systemctl start flagd
uninstall:
	systemctl disable flagd
	systemctl stop flagd
	rm /etc/systemd/system/flagd.service
	rm -f $(DESTDIR)$(PREFIX)/bin/flagd
lint:
	go install -v github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.1
	$(foreach module, $(ALL_GO_MOD_DIRS), ${GOPATH}/bin/golangci-lint run $(module)/...;)
lint-fix:
	go install -v github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.1
	$(foreach module, $(ALL_GO_MOD_DIRS), ${GOPATH}/bin/golangci-lint run --fix $(module)/...;)
install-mockgen:
	go install go.uber.org/mock/mockgen@v0.4.0
mockgen: install-mockgen
	cd core; mockgen -source=pkg/sync/http/http_sync.go -destination=pkg/sync/http/mock/http.go -package=syncmock
	cd core; mockgen -source=pkg/sync/grpc/grpc_sync.go -destination=pkg/sync/grpc/mock/grpc.go -package=grpcmock
	cd core; mockgen -source=pkg/sync/grpc/credentials/builder.go -destination=pkg/sync/grpc/credentials/mock/builder.go -package=credendialsmock
	cd core; mockgen -source=pkg/evaluator/ievaluator.go -destination=pkg/evaluator/mock/ievaluator.go -package=evalmock
	cd core; mockgen -source=pkg/sync/builder/syncbuilder.go -destination=pkg/sync/builder/mock/syncbuilder.go -package=middlewaremocksyncbuildermock
	cd flagd; mockgen -source=pkg/service/middleware/interface.go -destination=pkg/service/middleware/mock/interface.go -package=middlewaremock
generate-docs:
	cd flagd; go run ./cmd/doc/main.go

.PHONY: deploy-dev-env
export IMG?= ghcr.io/open-feature/flagd:latest
deploy-dev-env: undeploy-dev-env
	kubectl create ns "$(FLAGD_DEV_NAMESPACE)"
	envsubst '$${IMG}' < config/deployments/flagd/deployment.yaml | kubectl apply -f - -n "$(FLAGD_DEV_NAMESPACE)"
	kubectl apply -f config/deployments/flagd/service.yaml -n "$(FLAGD_DEV_NAMESPACE)"
	kubectl wait --for=condition=available deployment/flagd -n "$(FLAGD_DEV_NAMESPACE)" --timeout=300s

undeploy-dev-env:
	kubectl delete ns "$(FLAGD_DEV_NAMESPACE)" --ignore-not-found=true

run-zd-test:
	kubectl delete ns "$(ZD_TEST_NAMESPACE)" --ignore-not-found=true
	kubectl create ns "$(ZD_TEST_NAMESPACE)"
	ZD_TEST_NAMESPACE="$(ZD_TEST_NAMESPACE)" FLAGD_DEV_NAMESPACE=$(FLAGD_DEV_NAMESPACE) IMG="$(IMG)" IMG_ZD="$(IMG_ZD)" ./test/zero-downtime/zd_test.sh

run-flagd-proxy-zd-test:
	ZD_TEST_NAMESPACE_FLAGD_PROXY="$(ZD_TEST_NAMESPACE_FLAGD_PROXY)" FLAGD_PROXY_IMG="$(FLAGD_PROXY_IMG)" FLAGD_PROXY_IMG_ZD="$(FLAGD_PROXY_IMG_ZD)" ZD_CLIENT_IMG="$(ZD_CLIENT_IMG)" ./test/zero-downtime-flagd-proxy/zd_test.sh

# Markdown lint configuration
#
# - .markdownlintignore holds the configuration for files to be ignored
# - .markdownlint.yaml contains the rules for markdownfiles
MDL_DOCKER_VERSION := next
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
MDL_CMD := docker run -v $(ROOT_DIR):/workdir --rm 

.PHONY: markdownlint markdownlint-fix
markdownlint:
	$(MDL_CMD) davidanson/markdownlint-cli2:$(MDL_DOCKER_VERSION) "**/*.md" 

markdownlint-fix:
	$(MDL_CMD) davidanson/markdownlint-cli2:$(MDL_DOCKER_VERSION) --fix "**/*.md" 

.PHONY: pull-schemas-submodule
pull-schemas-submodule:
	git submodule update schemas

.PHONY: generate-proto-docs
generate-proto-docs: pull-schemas-submodule
	docker run --rm -v ${PWD}/$(DOCS_DIR)/reference/specifications:/out -v ${PWD}/schemas/protobuf:/protos pseudomuto/protoc-gen-doc --doc_opt=markdown,protos-with-toc.md flagd/evaluation/v1/evaluation.proto flagd/sync/v1/sync.proto \
	&& echo '<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->' > ${PWD}/$(DOCS_DIR)/reference/specifications/protos.md \
	&& sed '/^## Table of Contents/,/#top/d' ${PWD}/$(DOCS_DIR)/reference/specifications/protos-with-toc.md >> ${PWD}/$(DOCS_DIR)/reference/specifications/protos.md \
	&& rm -f ${PWD}/$(DOCS_DIR)/reference/specifications/protos-with-toc.md

# Update the schema at flagd.dev
# PUBLIC_JSON_SCHEMA_DIR above controls the dir (and therefore major version)
.PHONY: update-public-schema
update-public-schema: pull-schemas-submodule
	rm -f $(PUBLIC_JSON_SCHEMA_DIR)*.json
	cp schemas/json/*.json $(PUBLIC_JSON_SCHEMA_DIR)

.PHONY: run-web-docs
run-web-docs: generate-docs generate-proto-docs
	docker build -t flag-docs:latest .  --load \
	&& docker run --rm -it -p 8000:8000 -v ${PWD}:/docs flag-docs:latest

# Run the playground app in dev mode
# See the readme in the playground-app folder for more details
.PHONY: playground-dev
playground-dev:
	cd playground-app && npm ci && npm run dev

# Build the playground app
# See the readme in the playground-app folder for more details
.PHONY: playground-build
playground-build:
	cd playground-app && npm ci && npm run build

# Publish the playground app to the docs folder
# See the readme in the playground-app folder for more details
.PHONY: playground-publish
playground-publish: playground-build
	cp playground-app/dist/assets/index-*.js docs/playground/playground.js
