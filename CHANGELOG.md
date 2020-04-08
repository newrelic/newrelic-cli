<a name="unreleased"></a>
## [Unreleased]


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


[Unreleased]: https://github.com/newrelic/newrelic-client-go/compare/v0.6.2...HEAD
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
