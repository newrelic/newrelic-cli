[![Community Project header](https://github.com/newrelic/open-source-office/raw/master/examples/categories/images/Community_Project.png)](https://github.com/newrelic/open-source-office/blob/master/examples/categories/index.md#category-community-project)

# newrelic-cli

[![Testing](https://github.com/newrelic/newrelic-cli/workflows/Testing/badge.svg)](https://github.com/newrelic/newrelic-cli/actions)
[![Security Scan](https://github.com/newrelic/newrelic-cli/workflows/Security%20Scan/badge.svg)](https://github.com/newrelic/newrelic-cli/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/newrelic/newrelic-cli?style=flat-square)](https://goreportcard.com/report/github.com/newrelic/newrelic-cli)
[![GoDoc](https://godoc.org/github.com/newrelic/newrelic-cli?status.svg)](https://godoc.org/github.com/newrelic/newrelic-cli)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/newrelic/newrelic-cli/blob/main/LICENSE)
[![CLA assistant](https://cla-assistant.io/readme/badge/newrelic/newrelic-cli)](https://cla-assistant.io/newrelic/newrelic-cli)
[![Release](https://img.shields.io/github/v/release/newrelic/newrelic-cli?sort=semver)](https://github.com/newrelic/newrelic-cli/releases/latest)
[![Homebrew](https://img.shields.io/badge/dynamic/json.svg?url=https://formulae.brew.sh/api/formula/newrelic-cli.json&query=$.versions.stable&label=homebrew)](https://formulae.brew.sh/formula/newrelic-cli)

[![Docker Stars](https://img.shields.io/docker/stars/newrelic/cli.svg)](https://hub.docker.com/r/newrelic/cli)
[![Docker Pulls](https://img.shields.io/docker/pulls/newrelic/cli.svg)](https://hub.docker.com/r/newrelic/cli)
[![Docker Size](https://img.shields.io/docker/image-size/newrelic/cli.svg?sort=semver)](https://hub.docker.com/r/newrelic/k8s-operator)
[![Docker Version](https://img.shields.io/docker/v/newrelic/cli.svg?sort=semver)](https://hub.docker.com/r/newrelic/k8s-operator)

The New Relic CLI is an officially supported command line interface for New Relic, released as part of the [Developer Toolkit](https://newrelic.github.io/developer-toolkit/)

## Overview

The New Relic CLI is a project to consolidate some of the tools that New Relic
offers for managing resources. Current scope is limited while the framework is
being developed, but the tool as-is does perform a subset of tasks.

- Entity Search: Search for entities across all your New Relic accounts
- Entity Tagging: Manage tags across all of your entities
- Deployment Markers: Easily record an APM Application deployment within
  New Relic.

### Getting Started

For a quick guide on getting started with the New Relic CLI, see our [Getting
Started](https://github.com/newrelic/newrelic-cli/blob/main/docs/GETTING_STARTED.md)
page.

The latest New Relic CLI documentation is available [here](https://github.com/newrelic/newrelic-cli/blob/main/docs/cli/newrelic.md).

### Other Resources

There are a handful of other useful tools that this does not replace. Here are
some useful links to other tools that you might be interested in using at this
time.

- [NR1 CLI](https://developer.newrelic.com/build-tools/new-relic-one-applications/cli):
  Command line interface for managing development workflows for custom Nerdpacks on New Relic One.
- [New Relic Lambda CLI](https://github.com/newrelic/newrelic-lambda-cli): A
  CLI to install the New Relic AWS Lambda integration and layers.
- [New Relic Diagnostics](https://docs.newrelic.com/docs/agents/manage-apm-agents/troubleshooting/new-relic-diagnostics):
  A utility that automatically detects common problems with New Relic agents.

## Installation

Installation options are available for various platforms.

### MacOS

Install the New Relic CLI on MacOS via [`homebrew`](https://brew.sh/). With `homebrew` installed, run:

```
brew install newrelic-cli
```

### Windows

Installation is supported on 64-bit Windows.

#### Scoop

```powershell
scoop bucket add newrelic-cli https://github.com/newrelic/newrelic-cli.git
scoop install newrelic-cli
```

#### Chocolatey

```powershell
choco install newrelic-cli
```

#### Standalone installer

A standalone MSI installer is available on the GitHub releases page. You can download the installer for the latest version [here](https://github.com/newrelic/newrelic-cli/releases).

#### Powershell

Silent installation of the latest version of the CLI can be achieved via the follwing Powershell command:

```powershell
(New-Object System.Net.WebClient).DownloadFile("https://github.com/newrelic/newrelic-cli/releases/latest/download/NewRelicCLIInstaller.msi", "$env:TEMP\NewRelicCLIInstaller.msi"); `
msiexec.exe /qn /i "$env:TEMP\NewRelicCLIInstaller.msi" | Out-Null; `
```

### Linux

Linux binaries can be installed via [Snapcraft](https://snapcraft.io/). With the `snapd` daemon installed, run:

```
sudo snap install newrelic-cli
```

### Pre-built binaries

Pre-built binaries are created on the GitHub releases page for all of the above platforms. You can download the latest releases [here](https://github.com/newrelic/newrelic-cli/releases). The binaries and their checksums are signed and can be verified against the Developer Toolkit team's [public PGP key](https://newrelic.github.io/developer-toolkit/developer-toolkit.asc).

Verify that the fingerprint for the downloaded key matches the following:

```
gpg --fingerprint developer-toolkit-team@newrelic.com
86BE 01DA 9B1D A1FC F828  1409 DC9F C6B1 FCE4 7986
```

When verifying pre-built binaries and checksums, use the long format (the short format is not secure). For example:

```
gpg --keyid-format long --verify checksums.txt.sig checksums.txt
```

### Docker

There is an official [docker image](https://hub.docker.com/r/newrelic/cli) that can be utilized for running commands as well.

## Example Usage

#### Querying an APM application (using the Docker image)

```bash
# Pull the latest container
$ docker pull newrelic/cli

# Run the container interactively, remove it once the command exists
# Also must pass $NEW_RELIC_API_KEY to the container
$ docker run -it --rm \
    -e NEW_RELIC_API_KEY \
    newrelic/cli \
    apm application get --name WebPortal --accountId 2508259

[
  {
    "AccountID": 2508259,
    "ApplicationID": 204261368,
    "Domain": "APM",
    "EntityType": "APM_APPLICATION_ENTITY",
    "GUID": "MjUwODI1OXxBUE18QVBQTElDQVRJT058MjA0MjYxMzY4",
    "Name": "WebPortal",
    "Permalink": "https://one.newrelic.com/redirect/entity/MjUwODI1OXxBUE18QVBQTElDQVRJT058MjA0MjYxMzY4",
    "Reporting": true,
    "Type": "APPLICATION"
  }
]
```

See the [Getting Started guide](docs/GETTING_STARTED.md) for a more in-depth introduction to the capabilities of the New Relic CLI.

### Getting Help

In order to get help about what commands are available, the trusty `--help`
flag is here to assist. Alternatively, using just the `help` subcommand also works.

```
newrelic --help
newrelic help
```

Help is also available for the nested sub-commands. For example, the with the
following command, you can retrieve help for the `apm` sub-command.

```
newrelic apm --help
newrelic help apm
```

Using the CLI in this way, users are able to inspect what commands are
available, with some instruction on their usage.

### Patterns

Throughout the help, you may notice common patterns. The term `describe` is
used to perform list or get operations, while the `create` and `delete` terms
are used to construct or destroy an item, respectively.

## Development

### Requirements

- Go 1.16.0+
- GNU Make
- git

### Building

The `newrelic` command will be built in `bin/ARCH/newrelic`, where `ARCH` is either `linux`, `darwin`, or `windows`, depending on your build environment. You can run it directly from there or install it by moving it to a directory in your `PATH`.

```
# Default target is 'build'
$ make

# Explicitly run build
$ make build

# Locally test the CI build scripts
# make build-ci
```

### Testing

Before contributing, all linting and tests must pass. Tests can be run directly via:

```
# Tests and Linting
$ make test

# Only unit tests
$ make test-unit

# Only integration tests
$ make test-integration
```

### Working with recipes

#### Core recipe library

A core library of installation recipes is included with the CLI for use within the
`install` command. Recipe files are syndicated from [open-install-library](https://github.com/newrelic/open-install-library)
and embedded in the CLI binary at release time. To fetch the latest recipe library
while developing, the following make target can be used:

```
make recipes
```

Recipe files are stored in `internal/install/recipes/files`. Once files have been
fetched, they will be included in future CLI builds. If a particular version of
the recipe library is desired, the archive download URL can be passed to the make
target via the `RECIPES_ARCHIVE_URL` option:

```
make recipes RECIPES_ARCHIVE_URL=https://github.com/newrelic/open-install-library/releases/download/v0.50.0/recipes.zip
```

To clean recipe files, use the `recipes-clean` target:

```
make recipes-clean
```

#### Custom recipe files

A path can also be passed to the `--localRecipes` flag when running the `install`
command. This will bypass the methods described above and load files from the designated
path.

### Commit Messages

Using the following format for commit messages allows for auto-generation of
the [CHANGELOG](CHANGELOG.md):

#### Format:

`<type>(<scope>): <subject>`

| Type       | Description           | Change log? |
| ---------- | --------------------- | :---------: |
| `chore`    | Maintenance type work |     No      |
| `docs`     | Documentation Updates |     Yes     |
| `feat`     | New Features          |     Yes     |
| `fix`      | Bug Fixes             |     Yes     |
| `refactor` | Code Refactoring      |     No      |

#### Scope

This refers to what part of the code is the focus of the work. For example:

**General:**

- `build` - Work related to the build system (linting, makefiles, CI/CD, etc)
- `release` - Work related to cutting a new release

**Package Specific:**

- `newrelic` - Work related to the New Relic package
- `http` - Work related to the `internal/http` package
- `alerts` - Work related to the `pkg/alerts` package

### Documentation

**Note:** This requires the repo to be in your GOPATH [(godoc issue)](https://github.com/golang/go/issues/26827)

```
$ make docs
```

## Community Support

New Relic hosts and moderates an online forum where you can interact with New Relic employees as well as other customers to get help and share best practices.

- [Roadmap](https://newrelic.github.io/developer-toolkit/roadmap/) - As part of the Developer Toolkit, the roadmap for this project follows the same RFC process
- [Issues or Enhancement Requests](https://github.com/newrelic/newrelic-cli/issues) - Issues and enhancement requests can be submitted in the Issues tab of this repository. Please search for and review the existing open issues before submitting a new issue.
- [Contributors Guide](CONTRIBUTING.md) - Contributions are welcome (and if you submit a Enhancement Request, expect to be invited to contribute it yourself :grin:).
- [Community discussion board](https://discuss.newrelic.com/c/build-on-new-relic/developer-toolkit) - Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub.

Please do not report issues with the CLI to New Relic Global Technical Support. Instead, visit the [`Explorers Hub`](https://discuss.newrelic.com/c/build-on-new-relic) for troubleshooting and best-practices.

## Issues / Enhancement Requests

Issues and enhancement requests can be submitted in the [Issues tab of this repository](../../issues). Please search for and review the existing open issues before submitting a new issue.

## Contributing

Contributions are welcome (and if you submit a Enhancement Request, expect to be invited to contribute it yourself :grin:). Please review our [Contributors Guide](CONTRIBUTING.md).

Keep in mind that when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. If you'd like to execute our corporate CLA, or if you have any questions, please drop us an email at opensource@newrelic.com.

## Open Source License

This project is distributed under the [Apache 2 license](LICENSE).
