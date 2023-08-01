# flagd

Flagd is a feature flag daemon. It is a ready-made, open source, [OpenFeature](https://openfeature.dev) compliant feature flag backend system.

- OpenFeature compliant and speaks your language.
- Easy to extend to new languages.
- Supports multiple data sources simultaneously.
- Feature Flag updates occur in near real-time.
- Contains a powerful and flexible rule targeting engine and deterministic percentage-based rollouts.
- Flag evaluation statistics and metrics are exposed and compatible with Prometheus.

![flagd architecture](images/flagd-logical-architecture.jpg)

## flagd concepts

Whether you choose to run flagd on a Kubernetes cluster or outside of Kubernetes, there are concepts that apply to both equally.

Start your flagd learning journey here: [flagd concepts](concepts)

## Running flagd

### Kubernetes or Non-Kubernetes?

Now that you know the concepts, it's time to decide how to run flagd.

flagd can run on Kubernetes and in non-k8s environments.
Choose your operating environment to learn more:

- [flagd on Kubernetes](k8s)
- [flagd running outside of Kubernetes](nonk8s)
