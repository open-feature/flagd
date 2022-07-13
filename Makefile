IMG ?= flagd:latest
PHONY: .docker-build .build .run
PREFIX=/usr/local
guard-%:
	@ if [ "${${*}}" = "" ]; then \
        echo "Environment variable $* not set"; \
        exit 1; \
    fi
generate: guard-GOPATH
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	cd schemas/protobuf && ${GOPATH}/bin/buf generate && cd ../..
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
