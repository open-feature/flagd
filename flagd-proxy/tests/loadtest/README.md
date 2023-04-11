# flagd Proxy Profiling

This go module contains a profiling tool for the `flagd-proxy`. Starting `n` watchers against a single flag configuration resource to monitor the effects of server load and flag configuration definition size on the response time between a configuration change and all watchers receiving the configuration change.

## Pseudo Code

1. Parse configuration file referenced as the only startup argument
1. Loop for each defined repeat
1. Write to the target file using the start configuration
1. Start `n` watchers for the resource using a grpc sync definining the selector as `file:TARGET-FILE`
1. Wait for all watchers to receive their first configuration change event (which will contain the full configuration object)
1. Flush the change event channel to ensure there are no previous events
1. Trigger a configuration change event by writing the end configuration to the target file
8. Time how long it takes for all watchers to receive the new configuration

## Example

run the flagd-proxy locally (from the project root):

```
go run flagd-proxy/main.go start --port 8080
```

run the flagd-proxy-profiler (from the project root):

```
go run flagd-proxy/tests/loadtest/main.go ./flagd-proxy/tests/loadtest/config/config.json
```

Once the tests have been run the results can be found in ./flagd-proxy/tests/loadtest/profiling-results.json

## Sample Configuration

```
{
    "triggerType": "filepath",
    "fileTriggerConfig": {
        "startFile":"./start-spec.json",
        "endFile":"./config/end-spec.json",
        "targetFile":"./target.json"
    },
    "handlerConfig": {
        "filePath": "./target.json",
        "outFile":"./profiling-results.json",
        "host": "localhost",
        "port": 8080,
    },
    "tests": [
        {
            "watchers": 10000,
            "repeats": 5,
            "delay": 2000000000 
        }
    ]
}
```
