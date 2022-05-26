IMG=flagd:latest
docker-build:
	docker buildx build --platform="linux/amd64,linux/arm64" -t ${IMG} .
build:
	go build -o flagd
run:
	go run main.go start -f examples/example_flags.json
