---
description: flagd provider configuration
---

# flagd Provider Configuration

# 

Expose means to configure the provider aligned with the following priority system (highest to lowest).

```mermaid
flowchart LR
    explicit configuration -->|highest priority| environment-variables -->|lowest priority| defaults
```

### Explicit configuration

This takes the form of parameters to the provider's constructor, it has the highest priority.

### Environment variables

Read environment variables with sensible defaults (before applying the values explicitly declared to the constructor).

| Option name                 | Environment variable name             | Type    | Options      | Default                                |
| --------------------------- | ------------------------------------- | ------- | ------------ | -------------------------------------- |
| host                        | FLAGD_PROXY_HOST                      | string  |              | localhost                              |
| port                        | FLAGD_PROXY_PORT                      | number  |              | 8013                                   |
| tls                         | FLAGD_PROXY_TLS                       | boolean |              | false                                  |
| socketPath                  | FLAGD_PROXY_SOCKET_PATH               | string  |              |                                        |
| certPath                    | FLAGD_PROXY_SERVER_CERT_PATH          | string  |              |                                        |
| sourceURI                   | FLAGD_SOURCE_URI                      | string  |              |                                        |
| sourceProviderType          | FLAGD_SOURCE_PROVIDER_TYPE            | string  |              | grpc                                   |
| sourceSelector              | FLAGD_SOURCE_SELECTOR                 | string  |              |                                        |
| maxSyncRetries              | FLAGD_MAX_SYNC_RETRIES                | int     |              | 0 (0 means unlimited)                  |
| maxSyncRetryInterval        | FLAGD_MAX_SYNC_RETRY_INTERVAL         | int     |              | 60s                                    |