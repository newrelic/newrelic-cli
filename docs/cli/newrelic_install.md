## newrelic install

Install New Relic.

```
newrelic install [flags]
```

### Options

```
  -y, --assumeYes            use "yes" for all questions during install
      --debug                debug level logging
  -h, --help                 help for install
  -n, --recipe strings       the name of a recipe to install
  -c, --recipePath strings   the path to a recipe file to install
  -d, --skipDiscovery        skips discovery of recommended New Relic integrations
  -r, --skipIntegrations     skips installation of recommended New Relic integrations
  -l, --skipLoggingInstall   skips installation of New Relic Logging
  -t, --testMode             fakes operations for UX testing
      --trace                trace level logging
```

### Options inherited from parent commands

```
      --format string   output text format [YAML, JSON, Text] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic](newrelic.md)	 - The New Relic CLI

