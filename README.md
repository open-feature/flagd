# Flagd

Flagd is a simple command line tool for fetching and presenting feature flags to services. It is designed to conform to Open Feature schema for flag definitions.

<img src="images/of-flagd-0.png" width="560">      

## Example usage

Build the flagd binary:

```bash
make build
```

Start the process
```
./flagd start -f examples/example_flags.json --service-provider http --sync-provider filepath
```

This now provides an accessible http endpoint for the flags.
```
❯ curl localhost:8080
{ 
    "newWelcomeMessage": {
      "state": "disabled"
    },
    "hexColor": {
      "returnType": "string",
      "variants": {
        "red": "CC0000",
        "green": "00CC00",
```