# Integration tests

The continuous integration runs a set of [gherkin integration tests](https://github.com/open-feature/test-harness/blob/main/features).
If you'd like to run them locally, first pull the `test-harness` git submodule
```
git submodule update --init --recursive
```
then build the `flagd` binary
```
make build
```
then run the `flagd` binary
```
./flagd start -f file:test-harness/testing-flags.json
```
and finally run
```
make integration-test
```

## TLS

To run the integration tests against a `flagd` instance configured to use TLS, do the following:

Generate a cert and key in the repository root
```
openssl req -x509 -out localhost.crt -keyout localhost.key \
  -newkey rsa:2048 -nodes -sha256 \
  -subj '/CN=localhost' -extensions EXT -config <( \
   printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
```
build the `flagd` binary
```
make build
```
then run the `flagd` binary with tls configuration
```
./flagd start -f file:test-harness/testing-flags.json -c ./localhost.crt -k ./localhost.key
```
finally, either run the tests with an explicit path to the certificate:
```
make ARGS="-tls true -cert-path ./../../localhost.crt" integration-test
```
or, run without the path, defaulting to the host's root certificate authorities set (for this to work, the certificate must be registered and trusted in the host's system certificates)
```
make ARGS="-tls true" integration-test
```
