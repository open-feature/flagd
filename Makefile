IMG ?= flagd:latest
PHONY: .docker-build .build .run .mockgen
PREFIX=/usr/local
guard-%:
	@ if [ "${${*}}" = "" ]; then \
        echo "Environment variable $* not set"; \
        exit 1; \
    fi
generate:
	git submodule update --init --recursive
	cp schemas/json/flagd-definitions.json pkg/eval/flagd-definitions.json
	go install github.com/bufbuild/buf/cmd/buf@latest
	cd schemas/protobuf && buf generate --template buf.gen.go-server.yaml
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
run-ssl: generate local-certs
	go run main.go start -f config/samples/example_flags.json -c openfeature.crt -k openfeature.key
install:
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
install-mockgen:
	go install github.com/golang/mock/mockgen@v1.6.0
mockgen: install-mockgen
	mockgen -source=pkg/sync/http_sync.go -destination=pkg/sync/mock/http.go -package=syncmock
local-certs:
	openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
  -keyout openfeature.key -out openfeature.crt -subj "/CN=openfeature" \
  -addext "subjectAltName=DNS:openfeature.dev,DNS:www.openfeature.dev"