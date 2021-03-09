<a name="v0.19.1"></a>
## [v0.19.1] - 2021-03-05
### Bug Fixes
- Replace Detected observability gaps by Data Gaps

<a name="v0.19.0"></a>
## [v0.19.0] - 2021-03-02
### Bug Fixes
- **install:** allow skipDiscovery and skipLoggingInstall flags to work together

<a name="v0.18.32"></a>
## [v0.18.32] - 2021-02-25
<a name="v0.18.31"></a>
## [v0.18.31] - 2021-02-25
<a name="v0.18.30"></a>
## [v0.18.30] - 2021-02-24
<a name="v0.18.29"></a>
## [v0.18.29] - 2021-02-23
### Bug Fixes
- remove host filtering from InstallTarget
- **install:** tweak message
- **install:** fix lint
- **install:** Change fatal error message when failing to install

<a name="v0.18.28"></a>
## [v0.18.28] - 2021-02-19
### Bug Fixes
- update mock names
- simplify to use OpenInstallationPreInstallConfiguration from types package
- add PreInstall to the output Recipe

<a name="v0.18.27"></a>
## [v0.18.27] - 2021-02-19
<a name="v0.18.26"></a>
## [v0.18.26] - 2021-02-18
<a name="v0.18.25"></a>
## [v0.18.25] - 2021-02-17
### Bug Fixes
- use updated field

<a name="v0.18.24"></a>
## [v0.18.24] - 2021-02-12
### Bug Fixes
- **install:** Ensure error bubbles up when executing only 1 recipe and failing

<a name="v0.18.23"></a>
## [v0.18.23] - 2021-02-11
### Bug Fixes
- **install:** make ctrl-c exit reliably

### Features
- **install:** remove buffer of stdout/stderr when executing go-task
- **install:** allow recipe to display output messages

<a name="v0.18.22"></a>
## [v0.18.22] - 2021-02-05
### Bug Fixes
- **install:** suppress prompts during e2e tests

<a name="v0.18.21"></a>
## [v0.18.21] - 2021-02-05
### Bug Fixes
- **install:** prefer to prompt only when advanced mode is enabled
- **install:** include displayName in the recipeFile struct
- **install:** remove redundant [y/n] from prompts
- **install:** mask secret variables

### Features
- **decode:** NR CLI Translate base64 encoded urls command
- **decode:** NR CLI Translate base64 encoded urls command
- **decode:** NR CLI Translate base64 encoded urls command
- **decode:** NR CLI Translate base64 encoded urls command

<a name="v0.18.20"></a>
## [v0.18.20] - 2021-01-27
<a name="v0.18.19"></a>
## [v0.18.19] - 2021-01-27
### Bug Fixes
- **install:** ensure prompt response are handled correctly

<a name="v0.18.18"></a>
## [v0.18.18] - 2021-01-25
### Bug Fixes
- **install:** capture guid from validation when alternate format

<a name="v0.18.17"></a>
## [v0.18.17] - 2021-01-22
### Bug Fixes
- **install:** avoid duplicate results when fetching recommendations
- **install:** switch to confirmation for post-question visibility

<a name="v0.18.16"></a>
## [v0.18.16] - 2021-01-20
<a name="v0.18.15"></a>
## [v0.18.15] - 2021-01-20
<a name="v0.18.14"></a>
## [v0.18.14] - 2021-01-20
<a name="v0.18.13"></a>
## [v0.18.13] - 2021-01-20
### Bug Fixes
- **install:** avoid duplicate recipes from service results

<a name="v0.18.12"></a>
## [v0.18.12] - 2021-01-14
### Bug Fixes
- **install:** avoid newline in prompt message
- **install:** fetch and report logging before starting install
- **logging:** set client logger level

<a name="v0.18.11"></a>
## [v0.18.11] - 2021-01-14
### Bug Fixes
- **install:** report available recipes as soon as we know the list
- **install:** use the received name when fetching the recipe
- **install:** lint for else condition

<a name="v0.8.12"></a>
## [v0.8.12] - 2021-01-12
### Bug Fixes
- **install:** avoid newline in prompt message

<a name="v0.8.11"></a>
## [v0.8.11] - 2021-01-08
### Bug Fixes
- **install:** fix few lint issues
- **spinner:** drop duplicate spinner from output

<a name="v0.18.10"></a>
## [v0.18.10] - 2021-01-06
### Bug Fixes
- **install:** return error when default value is needed and not provided
- **install:** skip linting maligned struct
- **install:** set better default value when running automatic
- **install:** print newline after banner

### Features
- **install:** add -y flag

<a name="v0.18.9"></a>
## [v0.18.9] - 2020-12-29
### Bug Fixes
- **install:** avoid duplicate installs for logging and infra
- **install:** include displayName in request and recipe constructor
- **install:** avoid prompting when user has specified a named recipe
- **install:** capture task output and print only when debug logging
- **install:** avoid nil pointer and extra matches for service results
- **install:** update recipe spec to support displayName

<a name="v0.18.8"></a>
## [v0.18.8] - 2020-12-23
### Bug Fixes
- **install:** ensure secret input is hidden
- **install:** skip account-based link if default profile does not exist

<a name="v0.18.7"></a>
## [v0.18.7] - 2020-12-18
### Bug Fixes
- **install:** fixes for end to end flow
- **install:** tidy up the permissions on new files

<a name="v0.18.6"></a>
## [v0.18.6] - 2020-12-18
<a name="v0.18.5"></a>
## [v0.18.5] - 2020-12-17
### Bug Fixes
- **install:** replace package ID with default value
- **install:** create default log folder if not exists

<a name="v0.18.4"></a>
## [v0.18.4] - 2020-12-17
### Bug Fixes
- **build:** skip go generate as part of build process
- **install:** ignore region string case when checking profile
- **install:** reduce sudo requirement of install.sh
- **install:** detect and warn for empty NRQL validation
- **install:** let dead processes stay dead
- **install:** use string type for ID returned from the API

<a name="v0.18.3"></a>
## [v0.18.3] - 2020-12-11
<a name="v0.18.2"></a>
## [v0.18.2] - 2020-12-09
### Bug Fixes
- **install:** update logMatch type to list
- **install:** update logMatch type to list

<a name="v0.18.1"></a>
## [v0.18.1] - 2020-12-08
### Bug Fixes
- **install:** wire up all installContext fields

<a name="v0.18.0"></a>
## [v0.18.0] - 2020-12-08
### Features
- **install:** sketching out recipe validation
- **install:** fetch recipes from recipe service

<a name="v0.17.1"></a>
## [v0.17.1] - 2020-11-24
### Bug Fixes
- **diagnostics:** update download URL

<a name="v0.17.0"></a>
## [v0.17.0] - 2020-11-23
### Bug Fixes
- **install:** fix meltMatch struct to match spec

### Features
- **apiaccess:** add generated apiAccess commands (prerelease)
- **install:** prompt for variable input
- **install:** implement a mock server for process-based task selection
- **install:** wire up process discovery with cloned nri-process-discovery code

<a name="v0.16.0"></a>
## [v0.16.0] - 2020-11-04
### Bug Fixes
- **internal/diagnose:** download udpates via https!

### Features
- **internal/diagnose:** lint command; break out commands & helpers
- **internal/diagnose:** add minimal command line options
- **newrelic:** integrate with nrdiag (prototype)

<a name="v0.15.2"></a>
## [v0.15.2] - 2020-10-29
### Bug Fixes
- duplicitous task running
- **linting:** remove unused function

<a name="v0.15.1"></a>
## [v0.15.1] - 2020-10-28
<a name="v0.15.0"></a>
## [v0.15.0] - 2020-10-28
### Features
- **profiles:** create a profile automatically if it's possible

<a name="v0.14.1"></a>
## [v0.14.1] - 2020-10-15
### Bug Fixes
- **build:** update changelog action for improved standards

### Documentation Updates
- update changelog

<a name="v0.14.0"></a>
## [v0.14.0] - 2020-09-30
### Features
- **nerdgraph:** implement tutone-generated mutation command alertsPolicyCreate
- **release:** [#45](https://github.com/newrelic/newrelic-client-go/issues/45) add support for command chaining

<a name="v0.13.0"></a>
## [v0.13.0] - 2020-08-27
### Documentation Updates
- **readme:** include installation notes for Chocolatey users

<a name="v0.12.0"></a>
## [v0.12.0] - 2020-07-24
### Features
- **reporting:** add junit reporting

<a name="v0.11.0"></a>
## [v0.11.0] - 2020-07-24
### Features
- release edge command
- add a command for posting custom events

<a name="v0.10.0"></a>
## [v0.10.0] - 2020-07-10
### Bug Fixes
- **config:** remove extra comment
- **config:** compatible case with any loglevel value
- **config:** ensure compatible case with nr1
- **release:** need to use https URL for upstream homebrew-core
- **release:** generate correct sha256 for homebrew release
- **scoop:** fix bot email address

### Documentation Updates
- **README:** Correct Building section

### Features
- **release:** attempt to automate syncing our homebrew-core fork with upstream
- **release:** add step to update newrelic-forks/homebrew-core with latest from upstream homebrew-code
- **release:** update master branch with latest code from upstream

<a name="v0.9.0"></a>
## [v0.9.0] - 2020-06-16
### Bug Fixes
- **output:** Fix FormatText to do more than just tables

### Features
- **config:** Add config option to show Pre-Release Features (resolves [#274](https://github.com/newrelic/newrelic-client-go/issues/274))
- **edge:** mark as pre-release
- **edge:** add infinite tracing support
- **nrql:** Add NRQL Query and History commands
- **output:** Add text output formatter to general package

<a name="v0.8.5"></a>
## [v0.8.5] - 2020-05-27
### Bug Fixes
- **build:** Docker makefile was missing binary name
- **build:** ignore Scoop commits during commit linting
- **chocolatey:** fix typo in chocolatey verfication.txt

<a name="v0.8.4"></a>
## [v0.8.4] - 2020-05-11
### Bug Fixes
- **build:** Fix linting in Github actions

### Documentation Updates
- update community support information
- add the OSS category header
- **extensions:** add cli extension documentation

<a name="v0.8.3"></a>
## [v0.8.3] - 2020-04-24
### Bug Fixes
- **build:** Ignore scoop commit messages
- **release:** fix relative path in WiX project

<a name="v0.8.2"></a>
## [v0.8.2] - 2020-04-24
### Bug Fixes
- **release:** perform a stricter find when searching published assets

<a name="v0.8.1"></a>
## [v0.8.1] - 2020-04-24
### Bug Fixes
- **release:** use new token for publishing

<a name="v0.8.0"></a>
## [v0.8.0] - 2020-04-24
### Bug Fixes
- **chocolatey:** use copyright longer than four characters
- **chocolatey:** use better path for msi placement
- **chocolatey:** start packaging duing in main build process
- **chocolatey:** continue cleaning template files
- **chocolatey:** clean up comments from template files

### Documentation Updates
- update command examples in Getting Started guide to reflect recent updates
- **installation:** update installation guide with more options
- **readme:** fix typo in pgp key URL
- **releases:** add link to the DTK public PGP key

### Features
- **build:** provide installation via Scoop (Windows)
- **chocolatey:** begin chocolatey build
- **packaging:** include rpm and deb builds in goreleaser
- **release:** add code signing for artifacts

<a name="v0.7.0"></a>
## [v0.7.0] - 2020-04-20
### Bug Fixes
- **build:** Goreleaser was running `make clean` which broke things when run from `make release-publish`
- **ci:** wrap git config values in quotes
- **ci:** pass git config a global option
- **ci:** chmod +x the brew PR script
- **ci:** update the snap app name to match the binary
- **ci:** revert snapcraft binary name
- **ci:** upgrade snapcraft grade
- **ci:** add a step to install snapcraft
- **ci:** wire snapcraft token into publish step
- **ci:** wire docker creds into publish step
- **ci:** fix yaml indentation
- **ci:** update the snap name to match the binary

### Features
- **ci:** automate updating of homebrew formula
- **docs:** Custom release notes for goreleaser
- **install:** include install script
- **installer:** add code signing for Win installer
- **installers:** add a WiX-based MSI
- **newrelic-cli:** create a PS1 installer for Windows
- **newrelic-cli:** create a PS1 installer for Windows
- **output:** Enable format selection globally, also plain/pretty printing
- **output:** Support YAML output
- **output:** Output package for central output handling
- **snapcraft:** include goreleaser config for snaps

<a name="v0.6.2-test"></a>
## [v0.6.2-test] - 2020-04-09
<a name="v0.6.2"></a>
## [v0.6.2] - 2020-04-08
### Bug Fixes
- **region:** Add custom decoder for region for NR1 compatibility

<a name="v0.6.1"></a>
## [v0.6.1] - 2020-04-07
### Bug Fixes
- **newrelic:** Fix command name replacement on build
- **region:** Region parsing from config did not allow lowercase which is required for backwards compat

<a name="v0.6.0"></a>
## [v0.6.0] - 2020-04-03
### Features
- **nerdstorage:** add command for managing nerdstorage documents

<a name="v0.5.0"></a>
## [v0.5.0] - 2020-03-27
### Bug Fixes
- **credentials:** Change profiles => profile, remove => delete (with aliases)
- **documentation:** Unhide documentation generation command

### Features
- **docs:** Add auto-generated CLI documentation
- **docs:** Add cobra generated documentation command (hidden)
- **workloads:** add a command to duplicate workloads
- **workloads:** add a command to update workloads
- **workloads:** add a command to delete workloads
- **workloads:** add a command to create workloads
- **workloads:** add a command to list workloads
- **workloads:** add a command to get a workload

<a name="v0.4.1"></a>
## [v0.4.1] - 2020-03-11
### Bug Fixes
- **apm:** Fix apm command flag parsing for accountId, applicationId
- **apm/application:** Fix search params to accept accountId
- **credentials:** Lowercase region on storage for compatibility with nr1 cli

<a name="v0.4.0"></a>
## [v0.4.0] - 2020-03-10
### Bug Fixes
- **apm:** required params should result in help display
- **build:** Force tag fetch on CI checkout
- **lint:** skip spellcheck on the output/ directory
- **release:** include / in regex parising for commit messages

### Features
- **apm/deployment:** Add all user defined fields to the deployment creation
- **entities/search:** Return single object instead of array on single value result

<a name="v0.3.0"></a>
## [v0.3.0] - 2020-03-06
### Bug Fixes
- **newrelic:** Do not log fatal error if Cobra is printing out the help screen
- **newrelic:** avoid duplicate error message output

### Documentation Updates
- include information on getting started
- **newrelic:** update the help screens for consistency

### Features
- **entities:** add ability to map entity search result fields via flag
- **nerdgraph:** add nerdgraph command with query subcommand

<a name="v0.2.3"></a>
## [v0.2.3] - 2020-03-04
### Bug Fixes
- **build:** Allow overriding the version on make (needed for Homebrew local build)

<a name="v0.2.2"></a>
## [v0.2.2] - 2020-03-04
### Bug Fixes
- **build:** Enable remote docker for CircleCI
- **build:** Remove version.go generation from make release
- **build:** Add docker login to release-push process

<a name="v0.2.1"></a>
## [v0.2.1] - 2020-03-03
### Bug Fixes
- **docker:** Use Entrypoint so binary is assumed

<a name="v0.2.0"></a>
## [v0.2.0] - 2020-03-03
### Bug Fixes
- **client:** Fix ENV var prefix to be consistent with NR standards
- **config:** set user agent and service name, add version package
- **credentials:** proper handling when removing default profile
- **docs:** Fix release badge link

### Documentation Updates
- **command:** improve short help text

### Features
- **apm:** include get command for APM applications
- **docker:** Add docker image building / push to make system

<a name="v0.1.0"></a>
## v0.1.0 - 2020-02-27
### Bug Fixes
- load additional API key from environment
- Set correct module in go.mod
- **client:** resolve api key env var collision
- **config:** set defaults before validating config
- **config:** invert conditional for determining default fields
- **credentials:** allow setting profile if directory doesn't exist

### Documentation Updates
- Include overview documentation
- **entities:** include some examples and longer help

### Features
- **apm:** implement apm deployment marker retrieval
- **apm:** implement apm deployment create/delete
- **build:** Add docker handling to make system (build/clean/run)
- **build:** Create basic Dockerfile
- **completion:** include completion command for shell completion
- **config:** Add basic config loading
- **config:** write config file if none exists
- **config:** add remaining config methods
- **config:** add list method
- **config:** Add log level configuration
- **credentials:** implement initial credential management
- **credentials:** set default profile if adding one for the first time
- **credentials:** allow overriding api keys via env vars
- **entities:** add ability to filter entities search by entity type, tag, alert severity, domain, and reporting
- **entities:** add entity tag retrieval
- **entities:** implement entities tag and tag value deletion
- **entities:** implement add/replace tags
- **entities:** add entity search
- **profile:** Enable reading of profiles and use Region/APIKey from default profile
- **profile:** Add listing of profiles to command

[Unreleased]: https://github.com/newrelic/newrelic-client-go/compare/v0.19.1...HEAD
[v0.19.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.19.0...v0.19.1
[v0.19.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.32...v0.19.0
[v0.18.32]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.31...v0.18.32
[v0.18.31]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.30...v0.18.31
[v0.18.30]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.29...v0.18.30
[v0.18.29]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.28...v0.18.29
[v0.18.28]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.27...v0.18.28
[v0.18.27]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.26...v0.18.27
[v0.18.26]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.25...v0.18.26
[v0.18.25]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.24...v0.18.25
[v0.18.24]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.23...v0.18.24
[v0.18.23]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.22...v0.18.23
[v0.18.22]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.21...v0.18.22
[v0.18.21]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.20...v0.18.21
[v0.18.20]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.19...v0.18.20
[v0.18.19]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.18...v0.18.19
[v0.18.18]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.17...v0.18.18
[v0.18.17]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.16...v0.18.17
[v0.18.16]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.15...v0.18.16
[v0.18.15]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.14...v0.18.15
[v0.18.14]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.13...v0.18.14
[v0.18.13]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.12...v0.18.13
[v0.18.12]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.11...v0.18.12
[v0.18.11]: https://github.com/newrelic/newrelic-client-go/compare/v0.8.12...v0.18.11
[v0.8.12]: https://github.com/newrelic/newrelic-client-go/compare/v0.8.11...v0.8.12
[v0.8.11]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.10...v0.8.11
[v0.18.10]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.9...v0.18.10
[v0.18.9]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.8...v0.18.9
[v0.18.8]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.7...v0.18.8
[v0.18.7]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.6...v0.18.7
[v0.18.6]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.5...v0.18.6
[v0.18.5]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.4...v0.18.5
[v0.18.4]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.3...v0.18.4
[v0.18.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.2...v0.18.3
[v0.18.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.1...v0.18.2
[v0.18.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.0...v0.18.1
[v0.18.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.17.1...v0.18.0
[v0.17.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.17.0...v0.17.1
[v0.17.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.16.0...v0.17.0
[v0.16.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.15.2...v0.16.0
[v0.15.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.15.1...v0.15.2
[v0.15.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.15.0...v0.15.1
[v0.15.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.14.1...v0.15.0
[v0.14.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.14.0...v0.14.1
[v0.14.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.13.0...v0.14.0
[v0.13.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.12.0...v0.13.0
[v0.12.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.11.0...v0.12.0
[v0.11.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.10.0...v0.11.0
[v0.10.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.9.0...v0.10.0
[v0.9.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.8.5...v0.9.0
[v0.8.5]: https://github.com/newrelic/newrelic-client-go/compare/v0.8.4...v0.8.5
[v0.8.4]: https://github.com/newrelic/newrelic-client-go/compare/v0.8.3...v0.8.4
[v0.8.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.8.2...v0.8.3
[v0.8.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.8.1...v0.8.2
[v0.8.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.8.0...v0.8.1
[v0.8.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.7.0...v0.8.0
[v0.7.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.6.2-test...v0.7.0
[v0.6.2-test]: https://github.com/newrelic/newrelic-client-go/compare/v0.6.2...v0.6.2-test
[v0.6.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.6.1...v0.6.2
[v0.6.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.6.0...v0.6.1
[v0.6.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.5.0...v0.6.0
[v0.5.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.4.1...v0.5.0
[v0.4.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.4.0...v0.4.1
[v0.4.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.2.3...v0.3.0
[v0.2.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.2.2...v0.2.3
[v0.2.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.2.1...v0.2.2
[v0.2.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.2.0...v0.2.1
[v0.2.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.1.0...v0.2.0
