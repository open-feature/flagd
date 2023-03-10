Welcome!

There are a few things to consider before contributing to flagd.

Firstly, there's [a code of conduct](https://github.com/open-feature/.github/blob/main/CODE_OF_CONDUCT.md).
TLDR: be respectful.

Any contributions are expected to include unit tests. These can be validated with `make test` or the automated github workflow will run them on PR creation.

This project uses a go workspace, to setup the project run
```shell
make workspace-init
```

The go version in the `go.work` is the currently supported version of go.

The project uses remote buf packages, changing the remote generation source will require a one-time registry configuration:
```shell
export GOPRIVATE=buf.build/gen/go
```

Thanks! Issues and pull requests following these guidelines are welcome.
