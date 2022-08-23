IMG ?= flagd:latest
PHONY: .docker-build .build .run .mockgen
PREFIX=/usr/local

guard-%:
	@ if [ "${${*}}" = "" ]; then \
        echo "Environment variable $* not set"; \
        exit 1; \
    fi
docker-build:
	docker buildx build --build-arg=VERSION="$$(git describe --tags --abbrev=0)" --build-arg=COMMIT="$$(git rev-parse --short HEAD)" --build-arg DATE="$$(date +%FT%TZ)" --platform="linux/ppc64le,linux/s390x,linux/amd64,linux/arm64" -t ${IMG} .
docker-push:
	docker buildx build --push --build-arg=VERSION="$$(git describe --tags --abbrev=0)" --build-arg=COMMIT="$$(git rev-parse --short HEAD)" --build-arg DATE="$$(date +%FT%TZ)" --platform="linux/ppc64le,linux/s390x,linux/amd64,linux/arm64" -t ${IMG} .
build:
	go build -ldflags "-X main.version=dev -X main.commit=$$(git rev-parse --short HEAD) -X main.date=$$(date +%FT%TZ)" -o flagd
test:
	go test -cover ./...
run:
	go run main.go start -f config/samples/example_flags.json
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
