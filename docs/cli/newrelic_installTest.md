## newrelic installTest

Run a UX test of the install command.

```
newrelic installTest [flags]
```

### Options

```
  -y, --assumeYes             use "yes" for all questions during install
      --debug                 debug level logging
  -h, --help                  help for installTest
  -n, --recipe strings        the name of a recipe to install
  -c, --recipePath strings    the path to a recipe file to install
  -d, --skipDiscovery         skips discovery of recommended New Relic integrations
  -r, --skipIntegrations      skips installation of recommended New Relic integrations
  -l, --skipLoggingInstall    skips installation of New Relic Logging
  -s, --testScenario string   test scenario to run, defaults to BASIC.  Valid values are BASIC,LOG_MATCHES,FAIL,STITCHED_PATH,CANCELED,DISPLAY_EXPLORER_LINK (default "BASIC")
      --trace                 trace level logging
```

### Options inherited from parent commands

```
      --format string   output text format [JSON, Text, YAML] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic](newrelic.md)	 - The New Relic CLI

