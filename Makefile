IMG ?= flagd:latest
PHONY: .docker-build .build .run
PREFIX=/usr/local
guard-%:
	@ if [ "${${*}}" = "" ]; then \
        echo "Environment variable $* not set"; \
        exit 1; \
    fi
generate: guard-GOPATH
	git submodule update --init --recursive
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	cp schemas/json-schema/flagd-definitions.json pkg/eval/flagd-definitions.json
	${GOPATH}/bin/oapi-codegen --config=./config/open_api_gen_config.yml ./schemas/openapi/provider.yml
docker-build: generate
	docker buildx build --platform="linux/ppc64le,linux/s390x,linux/amd64,linux/arm64" -t ${IMG} .
docker-push: generate
	docker buildx build --push --platform="linux/ppc64le,linux/s390x,linux/amd64,linux/arm64" -t ${IMG} .
build: generate
	go build -o flagd
test: generate
	go test -cover ./...
run: generate
	go run main.go start -f config/samples/example_flags.json
install: build
	cp systemd/flagd.service /etc/systemd/system/flagd.service
	mkdir -p /etc/flagd
	cp systemd/flags.json /etc/flagd/flags.json
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp flagd $(DESTDIR)$(PREFIX)/bin/flagd
	systemctl start flagd
uninstall:
	systemctl disable flagd
	systemctl stop flagd
	rm /etc/systemd/system/flagd.service
	rm -f $(DESTDIR)$(PREFIX)/bin/flagd
lint:
	go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	${GOPATH}/bin/golangci-lint run --deadline=3m --timeout=3m ./... # Run linters