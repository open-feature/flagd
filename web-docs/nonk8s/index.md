# Install flagd

Install and run flagd almost anywhere outside of Kubernetes.

## Download flagd binary or container

flagd can be run as a standalone-binary or container.

Kubernetes-native? flagd can also be run [as part of the Kubernetes Operator](../k8s/index.md).

There are many ways to get started with flagd.
Choose the method that best serves your requirements to get started.

## Release binary

Download pre-built binaries from <https://github.com/open-feature/flagd/releases>

## Docker

```shell
docker pull ghcr.io/open-feature/flagd:latest
```

## Homebrew

```shell
brew install flagd
```

## Snap

[flagd is available on snapcraft](https://snapcraft.io/flagd):

```shell
sudo snap install flagd
```

## Go binary

1. Install Go 1.19 or above
1. Run `go install github.com/open-feature/flagd/flagd@latest`

## Systemd service

Documentation for installing flagd as a systemd service can be found [here](systemservice.md)

## flagd start

Start flagd:

```shell
flagd start \
  --port 8013 \
  --uri https://raw.githubusercontent.com/open-feature/flagd/main/samples/example_flags.flagd.json
```

Or use docker:

_Note - In Windows, use WSL system for both the file location and Docker runtime. Mixed file systems does not work and this is a [limitation of Docker](https://github.com/docker/for-win/issues/8479)._

```shell
docker run \
    --rm -it \
    --name flagd \
    -p 8013:8013 \
    ghcr.io/open-feature/flagd:latest start \
    --uri https://raw.githubusercontent.com/open-feature/flagd/main/samples/example_flags.flagd.json
```

If you wish, download the file locally to make changes:

```sh
wget https://raw.githubusercontent.com/open-feature/flagd/main/samples/example_flags.flagd.json
```

In local mode, run flagd like this:

```sh
flagd start \
    --port 8013 \
    --uri file:./example_flags.flagd.json
```

Or use docker ( _Note - In Windows, this requires WSL system for both the file location and Docker runtime_):

```sh
docker run \
    --rm -it \
    --name flagd \
    -p 8013:8013 \
    -v $(pwd):/etc/flagd \
    ghcr.io/open-feature/flagd:latest start \
    --uri file:./etc/flagd/example_flags.flagd.json
```

`--uri` can be a local file or any remote endpoint. Use `file:` prefix for local files. eg. `--uri file:/path/to/example_flags.flagd.json`. `gRPC` and `http` have their own requirements. More information can be found [here](https://github.com/open-feature/flagd/blob/main/docs/configuration/configuration.md#uri-patterns).

## Multiple flag sources and flag merging logic

Multiple `--uri` parameters can be specified. In other words, flagd can retrieve flags from multiple sources simultaneously.

See the [flag merging page](flagmerging.md) for more information.

## Perform flag evaluations

Flagd is now ready to perform flag evaluations over either `HTTP(s)` or `gRPC`. This example utilizes `HTTP` via `cURL`.

Retrieve a `String` value:

```sh
curl -X POST "http://localhost:8013/schema.v1.Service/ResolveString" \
    -d '{"flagKey":"myStringFlag","context":{}}' -H "Content-Type: application/json"
```

For Windows we recommend using a [WSL](https://learn.microsoft.com/en-us/windows/wsl/install) terminal.
Otherwise, use the following with `cmd`:

```sh
set json={"flagKey":"myStringFlag","context":{}}
curl -i -X POST -H "Content-Type: application/json" -d %json:"=\"% "localhost:8013/schema.v1.Service/ResolveString"
```

Result:

```json
{
    "value": "val1",
    "reason": "DEFAULT",
    "variant":"key1"
}
```

Updates to the underlying flag store (e.g. JSON file) are reflected by flagd in realtime. No restarts required.

flagd also supports boolean, integer, float and object flag types. Read more on the [evaluation examples page](https://github.com/open-feature/flagd/blob/main/docs/usage/evaluation_examples.md)

## Integrate your application

Now that flagd is running, it is time to integrate into your application. Do this by using [an OpenFeature provider in a language of your choice](https://github.com/open-feature/flagd/blob/main/docs/usage/flagd_providers.md).
