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
  -a, --skipApm              skips installation for APM
  -d, --skipDiscovery        skips discovery of recommended New Relic integrations
  -i, --skipInfra            skips installation for infrastructure agent (only for targeted install)
  -r, --skipIntegrations     skips installation of recommended New Relic integrations
  -l, --skipLoggingInstall   skips installation of New Relic Logging
  -t, --testMode             fakes operations for UX testing
      --trace                trace level logging
```

### Options inherited from parent commands

```
      --format string   output text format [JSON, Text, YAML] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic](newrelic.md)	 - The New Relic CLI

