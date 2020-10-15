# Integrating nrdiag

## Intro

The idea here is for the `newrelic` cli to manage downloading and updating the `nrdiag` binary, and to provide access to nrdiag's functionality by execing it. The current implementation is a working sketch but needs elaboration to be ready for production.

## Design

The `diag run` ensures the presence of the `nrdiag` binary, then runs it. The binary lives in ~/.newrelic/bin, which is created as necessary. If the binary is not present, `nrdiag_latest.zip` is downloaded from https://download.newrelic.com/nrdiag and expanded and the appropriate binary for the arch is put in place.

## Next steps

The following is a list of enhancements that have come up so far. I think most of them would be necessary for an MVP, but that is a matter for discussion.

* Security review - is this a reasonably safe approach?
* Update the `nrdiag` binary when necessary. nrdiag already has a function to check its version against the latest available. With a bit of tinkering this could be made machine-usable by the `newrelic` cli.
* Pass arbitrary arguments to the `nrdiag` binary. The command has too many options for it to feel reasonable to implement them in the `newrelic` cli and keep them in sync. (I tried setting `Args: cobra.ArbitraryArgs` on the `nrdiag run` subcommand, but this didn't work.) We need the user to be able to run, eg, `newrelic diag -t Java/Config/ValidateSettings -c config.yml` and have those arguments be passed on when we exec `nrdiag`. Currently this causes a usage error. 
* Optional: add a pass-through for `newrelic agent config lint`? The functionality my team recently added to nrdiag was meant primarily for use as a stand-alone linter. It might be nice to expose that functionality under the `agent` command as well as through `diag`.
* Optional: A subcommand to remove any downloaded files would probably be good, at least for potential use with Docker.

