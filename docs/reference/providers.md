# Providers

flagd was built from the ground up to be [Openfeature-compliant](../concepts/feature-flagging.md#openfeature-compliance).
To use it in your application, you must use the [OpenFeature SDK](https://openfeature.dev/docs/reference/technologies/) for your language, along with the associated OpenFeature _provider_.
For more information about Openfeature providers, see the [OpenFeature documentation](https://openfeature.dev/docs/reference/concepts/provider).

Providers for flagd come in two flavors: those that are built to communicate with a flagd instance (over HTTP or gRPC) and those that embed flagd's evaluation engine directly (note that some providers are capable of operating in either mode). For more information on how to deploy and use flagd, see [architecture](../architecture.md) and [installation](../installation.md).

For a catalog of available flagd providers, check out the [OpenFeature ecosystem](https://openfeature.dev/ecosystem?instant_search%5Bquery%5D=flagd&instant_search%5BrefinementList%5D%5Btype%5D%5B0%5D=Provider) page.

For information on implementing a flagd provider, see the specifications for [RPC](./specifications/rpc-providers.md) and [in-process](./specifications/in-process-providers.md) providers.
