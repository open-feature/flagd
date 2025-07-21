# Welcome

There are a few things to consider before contributing to flagd.

Firstly, there's [a code of conduct](https://github.com/open-feature/.github/blob/main/CODE_OF_CONDUCT.md).
TLDR: be respectful.

Any contributions are expected to include unit tests.
These can be validated with `make test` or the automated github workflow will run them on PR creation.

## Development

### Prerequisites

You'll need:

- Go
- make
- docker

You'll want:

- curl (for calling HTTP endpoints)
- [grpcurl](https://github.com/fullstorydev/grpcurl) (for making gRPC calls)
- jq (for pretty printing responses)

### Workspace Initialization

This project uses a go workspace, to setup the project run

```shell
make workspace-init
```

The go version in the `go.work` is the currently supported version of go.

The project uses remote buf packages, changing the remote generation source will require a one-time registry configuration:

```shell
export GOPRIVATE=buf.build/gen/go
```

### Manual testing

flagd has a number of interfaces (you can read more about than at [flagd.dev](https://flagd.dev/)) which can be used to evaluate flags, or deliver flag configurations so that they can be evaluated by _in-process_ providers.

You can manually test this functionality by starting flagd (from the flagd/ directory) with `go run main.go start -f file:../config/samples/example_flags.flagd.json`.

NOTE: you will need `go, curl`

#### Remote single flag evaluation via HTTP1.1/Connect

```sh
# evaluates a single boolean flag
curl -X POST -d '{"flagKey":"myBoolFlag","context":{}}' -H "Content-Type: application/json" "http://localhost:8013/flagd.evaluation.v1.Service/ResolveBoolean" | jq
```

#### Remote single flag evaluation via HTTP1.1/OFREP

```sh
# evaluates a single boolean flag
curl -X POST  -d '{"context":{}}' 'http://localhost:8016/ofrep/v1/evaluate/flags/myBoolFlag' | jq
```

#### Remote single flag evaluation via gRPC

```sh
# evaluates a single boolean flag
grpcurl -import-path schemas/protobuf/flagd/evaluation/v1/ -proto evaluation.proto -plaintext -d '{"flagKey":"myBoolFlag"}' localhost:8013 flagd.evaluation.v1.Service/ResolveBoolean | jq
```

#### Remote bulk evaluation via via HTTP1.1/OFREP

```sh
# evaluates flags in bulk
curl -X POST  -d '{"context":{}}' 'http://localhost:8016/ofrep/v1/evaluate/flags' | jq
```

#### Remote bulk evaluation via gRPC

```sh
# evaluates flags in bulk
grpcurl -import-path schemas/protobuf/flagd/evaluation/v1/ -proto evaluation.proto -plaintext -d '{}' localhost:8013 flagd.evaluation.v1.Service/ResolveAll | jq
```

#### Flag configuration fetch via gRPC

```sh
# sends back a representation of all flags
grpcurl -import-path schemas/protobuf/flagd/sync/v1/ -proto sync.proto -plaintext localhost:8015 flagd.sync.v1.FlagSyncService/FetchAllFlags | jq
```

#### Flag synchronization stream via gRPC

```sh
# will open a persistent stream which sends flag changes when the watched source is modified 
grpcurl -import-path schemas/protobuf/flagd/sync/v1/ -proto sync.proto -plaintext localhost:8015 flagd.sync.v1.FlagSyncService/SyncFlags | jq
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

   ```shell
   Signed-off-by: Joe Smith <joe.smith@example.com>
   ```

   You MUST use your legal name -- handles or other pseudonyms are not
   permitted.

   While you could manually add DCO sign-off to every commit, there is an easier
   way:

    1. Configure your git client appropriately. This is one-time setup.

       ```shell
       git config user.name <legal name>
       git config user.email <email address you use for GitHub>
       ```

       If you work on multiple projects that require a DCO sign-off, you can
       configure your git client to use these settings globally instead of only
       for Brigade:

       ```shell
       git config --global user.name <legal name>
       git config --global user.email <email address you use for GitHub>
       ```

    1. Use the `--signoff` or `-s` (lowercase) flag when making each commit.
       For example:

       ```shell
       git commit --message "<commit message>" --signoff
       ```

       If you ever make a commit and forget to use the `--signoff` flag, you
       can amend your commit with this information before pushing:

       ```shell
       git commit --amend --signoff
       ```

    1. You can verify the above worked as expected using `git log`. Your latest
       commit should look similar to this one:

       ```shell
       Author: Joe Smith <joe.smith@example.com>
       Date:   Thu Feb 2 11:41:15 2018 -0800

           Update README

           Signed-off-by: Joe Smith <joe.smith@example.com>
       ```

       Notice the `Author` and `Signed-off-by` lines match. If they do not, the
       PR will be rejected by the automated DCO check.

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
