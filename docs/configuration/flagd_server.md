## flagd server

Start flagd as a server

```
flagd server [flags]
```

### Options

```
  -p, --address string     Path this server binds to (default "localhost:9090")
  -c, --cert-path string   TLS certificate path
  -h, --help               help for server
  -k, --key-path string    TLS key path of the certificate
  -s, --secure             Start secure server
  -f, --source string      CRD with feature flag configurations
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.agent.yaml)
  -x, --debug           verbose logging
```

### SEE ALSO

* [flagd](flagd.md)	 - Flagd is a simple command line tool for fetching and presenting feature flags to services. It is designed to conform to Open Feature schema for flag definitions.

