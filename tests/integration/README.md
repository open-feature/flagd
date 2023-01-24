#### Integration tests

The continuous integration runs a set of [gherkin integration tests](https://github.com/open-feature/test-harness/blob/main/features).
If you'd like to run them locally, first pull the `test-harness` git submodule
```
git submodule update --init --recursive
```
then build your flagd image
```
make docker-build
```
then run your built flagd image
```
docker run -p 8013:8013 -v $PWD/test-harness/testing-flags.json:/testing-flags.json ghcr.io/open-feature/flagd:latest start -f file:/testing-flags.json
```
and finally run
```
make integration-test
```

Note: Testing against the flagd binary directly (rather than the docker image) introduces test flakiness.
