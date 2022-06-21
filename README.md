# Flagd

![build](https://img.shields.io/github/workflow/status/open-feature/flagd/ci)
![goversion](https://img.shields.io/github/go-mod/go-version/open-feature/flagd/main)
![version](https://img.shields.io/badge/version-pre--alpha-green)
![status](https://img.shields.io/badge/status-not--for--production-red)

Flagd is a simple command line tool for fetching and presenting feature flags to services. It is designed to conform to OpenFeature schema for flag definitions.

<img src="images/of-flagd-0.png" width="560">

## Example usage

1. Generate the prerequisites `make generate`
2. Build the flagd binary: `make build`
3. Start the process: `./flagd start -f config/samples/example_flags.json --service-provider http --sync-provider filepath`

This now provides an accessible http endpoint for the flags:

```
$ curl -X POST "localhost:8080/flags/myBoolFlag/resolve/boolean?default-value=true"
// {"reason":"STATIC","value":true}

$ curl -X POST "localhost:8080/flags/myStringFlag/resolve/string?default-value=hi"
// {"reason":"STATIC","value":"red"}

$ curl -X POST "localhost:8080/flags/myNumberFlag/resolve/number?default-value=13"
// {"reason":"STATIC","value":1}

$ curl -X POST "localhost:8080/flags/myObjectFlag/resolve/object?default-value=foo,bar"
// {"reason":"STATIC","value":{"color":"blue"}}

$ curl -X POST localhost:8080/flags/isColorYellow/resolve/boolean?default-value=true \
-d '{"color": "yellow"}'
// {"reason":"TARGETING_MATCH","value":true}
```

### Installation

#### Systemd

To install as a systemd service run `sudo make install` this will place the binary by default in `/usr/local/bin`

There will also be a default provider and sync enabled ( http / filepath ) both of which can be modified in the flagd.service.

Validation can be run with `systemctl status flagd`
And result similar to below will be seen

```
● flagd.service - "A generic feature flag daemon"
     Loaded: loaded (/etc/systemd/system/flagd.service; disabled; vendor preset: enabled)
     Active: active (running) since Mon 2022-05-30 12:19:55 BST; 5min ago
   Main PID: 64610 (flagd)
      Tasks: 7 (limit: 4572)
     Memory: 1.4M
     CGroup: /system.slice/flagd.service
             └─64610 /usr/local/bin/flagd start -f=/etc/flagd/flags.json

May 30 12:19:55 foo systemd[1]: Started "A generic feature flag daemon".
```

### Running in a container

1. `IMG=flagd-local make docker-build`
2. `docker run -p 8080:8080 -it flagd-local start --uri ./examples/example_flags.json`
