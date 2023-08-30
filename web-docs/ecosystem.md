# The flagd ecosystem

flagd relies upon many projects and tools.

We thank the following projects and always contribute upstream whenever and wherever it makes sense.

## Kubernetes and the Operator Pattern

flagd can be leveraged in Kubernetes clusters by installing the [OpenFeature Operator](k8s/index.md).

The OpenFeature Operator is, as the name suggests, written to follow the [Kubernetes Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

## ArtifactHub

The OpenFeature Operator is [listed on ArtifactHub](https://artifacthub.io/packages/helm/open-feature-operator/open-feature-operator).

## Helm

[Helm](https://helm.sh) is used to package and install the OpenFeature Operator.

## cert-manager

[cert-manager](https://cert-manager.io/) is a prerequisite of installing the OpenFeature Operator.

The OpenFeature Operator is a server that communicates with Kubernetes components within a cluster.
As such, it requires a means of authorizing requests between peers.
Cert manager handles authorization by adding certificates and certificate issuers as resource types in Kubernetes clusters.
This simplifies the process of obtaining, renewing, and using those certificates.

## gRPC

flagd offers [gRPC](https://grpc.io) in two ways:

- An interface between the OpenFeature SDK (ie. clients who want to query flagd).
- An interface between flagd and feature flag providers (via the SyncProvider concept).

## MurmurHash

[MurmurHash](https://github.com/aappleby/smhasher) (specificially [MurmurHash3](https://github.com/aappleby/smhasher/blob/master/src/MurmurHash3.cpp)) is used by flagd during [fractional evaluation](concepts/index.md#fractional-evaluation).
