# Welcome

There are a few things to consider before contributing to flagd.

Firstly, there's [a code of conduct](https://github.com/open-feature/.github/blob/main/CODE_OF_CONDUCT.md).
TLDR: be respectful.

Any contributions are expected to include unit tests.
These can be validated with `make test` or the automated github workflow will run them on PR creation.

This project uses a go workspace, to setup the project run

```shell
make workspace-init
```

The go version in the `go.work` is the currently supported version of go.

The project uses remote buf packages, changing the remote generation source will require a one-time registry configuration:

```shell
export GOPRIVATE=buf.build/gen/go
```

## DCO Sign-Off

A DCO (Developer Certificate of Origin) sign-off is a line placed at the end of
a commit message containing a contributor's "signature." In adding this, the
contributor certifies that they have the right to contribute the material in
question.

Here are the steps to sign your work:

1. Verify the contribution in your commit complies with the
   [terms of the DCO](https://developercertificate.org/).

1. Add a line like the following to your commit message:


## Conventional PR Titles

When raising PRs, please format according to [conventional commit standards](https://www.conventionalcommits.org/en/v1.0.0/#summary)

For example: `docs: some PR title here...`

Thanks!
Issues and pull requests following these guidelines are welcome.

## Markdown Lint and Markdown Lint Fix

PRs are expected to conform to markdown lint rules.

Therefore, run `make markdownlint-fix` to auto-fix _most_ issues.
Then commit the results.

For those issues that cannot be auto-fixed, run `make markdownlint`
then manually fix whatever it warns about.
