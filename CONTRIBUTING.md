Welcome!

There are a few things to consider before contributing to flagd.

Firstly, there's [a code of conduct](https://github.com/open-feature/.github/blob/main/CODE_OF_CONDUCT.md).
TLDR: be respectful.

Any contributions are expected to include unit tests. These can be validated with `make test` or the automated github workflow will run them on PR creation.

The go version in the `go.mod` is the currently supported version of go.

The project uses remote buf packages which will require a one-time registry configuration for local development:
```shell
export GOPRIVATE=buf.build/gen/go
```

Thanks! Issues and pull requests following these guidelines are welcome.
