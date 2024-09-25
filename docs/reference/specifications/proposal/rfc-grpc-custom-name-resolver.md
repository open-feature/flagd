# gRPC Custom Name Resolver Proposal (DRAFT)

## Details

|                        |                                    |
|------------------------|------------------------------------|
| **Feature Name**       | gRPC custom name resolver          |
| **Type**               | enhancement                        |
| **Related components** | gRPC source resolution             |

## Summary

gRPC by default supports DNS resolution which is currently being used e.g. "localhost:8013" in both
[core](https://github.com/open-feature/flagd/blob/main/core/pkg/sync/grpc/grpc_sync.go#L72-L74) and
providers e.g. [java](https://github.com/open-feature/java-sdk-contrib/blob/main/providers/flagd/src/main/java/dev/openfeature/contrib/providers/flagd/resolver/common/ChannelBuilder.java#L53-L55).
This covers most deployments, but with increased adoption of microservice-architecture, service discovery,
policy-enabled service meshes (e.g. istio, envoy, consul, etc) it's necessary to support custom routing and name resolution.

For such cases the gRPC core libraries support few alternative resolver* also expose the required interfaces to build custom implementations:

### Reference

* [Custom Name Resolution](https://grpc.io/docs/guides/custom-name-resolution/)
* [Java Client](https://grpc.github.io/grpc-java/javadoc/io/grpc/ManagedChannelBuilder.html#forTarget(java.lang.String))
* [Golang](https://pkg.go.dev/google.golang.org/grpc#NewClient)

A sample workflow of a deployment with proxy sidecar using `in-process` mode (for `rpc` the process is same)

```mermaid
sequenceDiagram
    participant flagd-sync.service
    participant proxy-sidecar-agent
    participant flagd-provider
    participant application


    flagd-provider-->proxy-sidecar-agent: Route traffic to backend sync service based on policy
    Note left of flagd-provider: In most cases the connection made to <br> `localhost` data plane port <br> e.g. localhost:9211

    flagd-provider-->flagd-sync.service: Get realtime flag update over gRPC stream
    Note left of proxy-sidecar-agent: Based on the host header i.e. authority <br> the tarrfic was then routed <br> backend service.
    loop
        flagd-provider->>flagd-provider: in-memory cache
    end
    application-->>flagd-provider: Check the state of a feature flag
    flagd-provider-->>application: Get the feature flag from in-memory cache <br/> run the evaluation logic and return final state
```

**Note:** There is small variation in supported alternative resolver e.g. java support `zooKeeper`

## Motivation

The main motivation is to support complex deployments with a generic custom name resolver using the interface
provided by gRPC core*.

**Note**: As of now only `java` and `golang` has the required interface to create custom resolver

## Detailed design

The idea is to

- allow a new config option to pass the [target](https://grpc.io/docs/guides/custom-name-resolution/#life-of-a-target-string) string
- reduce need to create/override existing implementations to simplify use of name-resolver

### Target String Pattern*

There is no restriction on naming but the string most comply with below standard

```text
scheme://authority/endpoint_name
```

#### Scheme Suggestions

We can choose our own scheme considering it's not used in gRPC core

| Scheme             | Comment | Status             |
|--------------------|---------|--------------------|
| **flagd://**       |         | Approved / Reject  |
| **grpc-remote://** |         | Approved / Reject  |
| **tbd**            |         | Approved / Reject  |

##### Authority

Authority needs to be a valid tcp endpoint of proxy/service discovery agent (passed by the user)
e.g. `localhost:9211`

##### Endpoint Name

The endpoint also specific to user deployment environment which define the `flagd` or `sync` backend
service name i.e. VirtualHost. This is used by the provided authority i.e. proxy service where to
route the traffic e.g.

```shell
$ grpcurl -vv -plaintext -authority flagd-sync.service 127.0.0.1:9211 list flagd.sync.v1.FlagSyncService
flagd.sync.v1.FlagSyncService.FetchAllFlags
flagd.sync.v1.FlagSyncService.GetMetadata
flagd.sync.v1.FlagSyncService.SyncFlags
```

##### Samples

* ``[ flagd || grpc-remote ]://[ 127.0.0.1:9211 ]/[ flagd-sync.service ]``
* ``[ flagd || grpc-remote ]://[ proxy.domain:443 ]/[ flagd-sync.service ]``

##### Drawbacks

* One of the big drawback was limited support of the language only `java` and `golang`
* Will introduce inconsistent user experience
* Will open the door for different use cases although this can be fixed by
providing sdks similar to [custom connector](https://github.com/open-feature/java-sdk-contrib/tree/main/providers/flagd#custom-connector)
* ...

## Alternatives

### Option-1

Allow users to override default ``authority`` header as shown above in `grpcurl`, the override option was
already supported by all major languages*

* [Golang](https://pkg.go.dev/google.golang.org/grpc#WithAuthority)
* [JAVA](https://grpc.github.io/grpc-java/javadoc/io/grpc/ForwardingChannelBuilder2.html#overrideAuthority(java.lang.String))
* [Python](https://grpc.github.io/grpc/python/glossary.html#term-channel_arguments)

this option is simple and easy to implement, although it will not cover all the cases it will at least help with proxy
setup where `host_header` was used to route traffic.

**Ref**:

Java PR: <https://github.com/open-feature/java-sdk-contrib/pull/949>

**Note**: JS, .NET, PHP still need to be explored if this options available

### Option-2

Only support the [xDS](https://grpc.io/docs/guides/custom-load-balancing/#service-mesh) protocol which already supported by gRPC core and doesn't require any custom
name resolver we can simply use any `target` string with `xds://` scheme. The big benefit of this approach was
it's going to be new stranded when it comes gRPC with service mesh and eliminate any custom implementation in `flagd`
and the gRPC core team actively adding more features e.g. mTLS

For more details refer the below document

* [gRPC xDS Feature](https://grpc.github.io/grpc/core/md_doc_grpc_xds_features.html)
* [gRPC xDS RFC](https://github.com/grpc/proposal/blob/master/A52-xds-custom-lb-policies.md)

### Option-3

TBD

## Unresolved questions

* What to do with un-supported languages
* Coming up with generic name resolver which will cover most of the cases not just proxy
* ....
