# flagd

## Quick Start

```shell
# Start flagd
docker run \
    --rm -it \
    --name flagd \
    -p 8013:8013 \
    ghcr.io/open-feature/flagd:latest start \
    --uri https://raw.githubusercontent.com/open-feature/flagd/main/samples/example_flags.flagd.json

# Query a feature flag called "isColorYellow"
# variant == "off", value == false
curl -X POST "http://localhost:8013/schema.v1.Service/ResolveBoolean" \
  -d '{"flagKey":"isColorYellow","context":{}}' \
  -H "Content-Type: application/json"
# {"reason":"TARGETING_MATCH","variant":"off"}

# Query "isColorYellow" flag again
# with runtime contextual information
# variant == "on", value == true
curl -X POST "http://localhost:8013/schema.v1.Service/ResolveBoolean" \
  -d '{"flagKey":"isColorYellow","context":{"color": "yellow"}}' \
  -H "Content-Type: application/json"
# {"value":true,"reason":"TARGETING_MATCH","variant":"on"}
```

flagd can also run:

- [As a standalone binary](https://flagd.dev/nonk8s/#release-binary)
- [As a Kubernetes Operator](https://flagd.dev/k8s/)
- [Installed using homebrew](https://flagd.dev/nonk8s/#homebrew)
- [Installed using snap](https://flagd.dev/nonk8s/#snap)

## What is flagd?

Flagd is a feature flag daemon. It is a ready-made, open source, [OpenFeature](https://openfeature.dev) compliant feature flag backend system.

- OpenFeature compliant and speaks your language.
- Easy to extend to new languages.
- Supports multiple data sources simultaneously.
- Feature Flag updates occur in near real-time.
- Contains a powerful and flexible rule targeting engine and deterministic percentage-based rollouts.
- Flag evaluation statistics and metrics are exposed and compatible with OpenTelemetry.

![flagd architecture](images/flagd-logical-architecture.jpg)

## flagd concepts

Whether you choose to run flagd on a Kubernetes cluster or outside of Kubernetes, there are concepts that apply to both equally.

Start your flagd learning journey here: [flagd concepts](concepts/index.md)

## Running flagd

### Kubernetes or Non-Kubernetes?

Now that you know the concepts, it's time to decide how to run flagd.

flagd can run on Kubernetes and in non-k8s environments.
Choose your mode of deployment to learn more:

- [flagd on Kubernetes](k8s/index.md)
- [flagd running outside of Kubernetes](nonk8s/index.md)
