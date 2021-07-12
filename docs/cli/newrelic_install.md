## newrelic install

Install New Relic.

```
newrelic install [flags]
```

### Options

```
  -y, --assumeYes             use "yes" for all questions during install
  -h, --help                  help for install
      --localRecipes string   a path to local recipes to load instead of service other fetching
  -n, --recipe strings        the name of a recipe to install
  -c, --recipePath strings    the path to a recipe file to install
  -s, --skipApm               skips installation for APM
  -i, --skipInfra             skips installation for infrastructure agent (only for targeted install)
  -r, --skipIntegrations      skips installation of recommended New Relic integrations
  -l, --skipLoggingInstall    skips installation of New Relic Logging
  -t, --testMode              fakes operations for UX testing
```

### Options inherited from parent commands

```
  -a, --accountId int    trace level logging
      --debug            debug level logging
      --format string    output text format [JSON, Text, YAML] (default "JSON")
      --plain            output compact text
      --profile string   the authentication profile to use
      --trace            trace level logging
```

### SEE ALSO

* [newrelic](newrelic.md)	 - The New Relic CLI

