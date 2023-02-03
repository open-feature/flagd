#### Integration tests

The continuous integration runs a set of [gherkin integration tests](https://github.com/open-feature/test-harness/blob/main/features).
If you'd like to run them locally, first pull the `test-harness` git submodule
```
git submodule update --init --recursive
```
then build the flagd binary
```
make build
```
then run the flagd binary
```
./flagd start -f file:test-harness/testing-flags.json
```
and finally run
```
make integration-test
```
