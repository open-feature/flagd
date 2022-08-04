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
	rm *.pem ||

	# 1. Generate CA's private key and self-signed certificate
	openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=FR/ST=Occitanie/L=Toulouse/O=Tech School/OU=Education/CN=*.techschool.guru/emailAddress=techschool.guru@gmail.com"

	echo "CA's self-signed certificate"
	openssl x509 -in ca-cert.pem -noout -text

	# 2. Generate web server's private key and certificate signing request (CSR)
	openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=FR/ST=Ile de France/L=Paris/O=PC Book/OU=Computer/CN=*.pcbook.com/emailAddress=pcbook@gmail.com"

	# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
	openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

	echo "Server's signed certificate"
	openssl x509 -in server-cert.pem -noout -text