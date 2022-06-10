IMG=flagd:latest
PHONY: .docker-build .build .run
PREFIX=/usr/local
generate:
	${GOPATH}/bin/oapi-codegen --config=./config/open_api_gen_config.yml ./schemas/openapi/provider.yml
docker-build:
	docker buildx build --platform="linux/ppc64le,linux/s390x,linux/amd64,linux/arm64" -t ${IMG} .
build: generate
	go build -o flagd
run:
	go run main.go start -f examples/example_flags.json
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
