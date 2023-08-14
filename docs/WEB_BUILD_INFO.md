# How https://flagd.dev is built

The flagd project has a website at https://flagd.dev

This guide is intended so that potential contributors can understand how, where and when the website comes into being.

## Website Structure

The website is built on [mkdocs](https://www.mkdocs.org/) and the [material theme](https://squidfunk.github.io/mkdocs-material/).

The [mkdocs.yml](../mkdocs.yml) file describes the high level metadata in addition to the main menu structure.

The website is generated from files in the [web-docs](../web-docs/) folder.

The website homepage is powered by the [index.md](../web-docs/index.md) file at the root of `/web-docs`

## What happens during a Website-Related PR?

Imagine you notice a typo on https://flagd.dev.

You find the page inside [/web-docs](../web-docs/) and fix the typo and commit the change (remembering to sign your commit).

You raise a PR with your change.

What next? What checks run against your change?

1. All commits in this PR are checked for being signed with DCO.
   Hint: Sign **all** commits with `-s`: `git commit -sm "my commit message"`
   If the DCO check fails, this bot will provide the copy / paste commands to fix.
   Click on the Details hyperlink in the PR and scroll to the bottom of the resulting page.
1. PR titles must be semantically correct.
   [The lint-pr GitHub action](../.github/workflows/lint-pr.yaml) checks that and warns if your PR title doesn't match requirements.
1. The markdown files are linted in [the markdown-checks GitHub Action](../.github/workflows/markdown-checks.yaml) which calls [make markdownlint](../Makefile#L103).
   Your markdown syntax [must follow a set of rules](https://github.com/DavidAnson/markdownlint/blob/main/doc/Rules.md).
   Hint: Run `make markdownlint-fix` and / or `make markdownlint` before raising the PR!
1. [The build GitHub Action](../.github/workflows/build.yaml) runs a job called [lint](../.github/workflows/build.yaml#L23)
   which in turn runs [make workspace-init](../Makefile#L12) and then [make lint](../Makefile#L62).
   TODO: @thomaspoignant to advise on what this does and if it's still relevant for documentation.
   That linting _appears_ to me to be related to Go code, not docs.
   Perhaps it can be re-scoped to run only for Go code?
1. [The build GitHub Action](../.github/workflows/build.yaml) runs a job called [docs-check](../.github/workflows/build.yaml#L38)
   which in turn runs [make workspace-init](../Makefile#L12) (duplication of above!).
   Then [make generate-docs](../Makefile#73) - the purpose of which is unclear.
1. [The build GitHub Action](../.github/workflows/build.yaml) runs a job called [test](../.github/workflows/build.yaml#L38) which checks out the repository, sets up Go, runs [make workspace-init](../Makefile#L12) (for this third time during this PR).
   [make test](../Makefile#L39) which defaults to [make test-core](../Makefile#L42) is run, which, as far as I can see, does nothing related to documentation.
   TODO: Re-scope or remove?

1. [The build GitHub Action](../.github/workflows/build.yaml) runs a job called [docker-local](../.github/workflows/build.yaml#L70) which, as far as I can see, does nothing at all related to documentation.
  Instead this job appears to build the Go code, scan it for vulnerabilities with Trivy and upload the results.
  TODO: Rebuild / rescope to remove from documentation PRs?

1. [The build GitHub Action](../.github/workflows/build.yaml) runs a job called [integration-test](../.github/workflows/build.yaml#L115).
   This job initialises a Go workspace (for the fourth time in this PR), builds the code, builds a flagd binary and runs integration tests.
   All of which appear irrelevant to the documentation.

   TODO: Rebuild / rescope to remove from documentation PRs?
1. A preview version of the site is built by Netflify.
1. The Netlify preview lint is added to the PR.
1. Netlify will rebuild the site for every new commit.


## How does the site get built and where is it deployed?

@beeme1mr to help author this section.

