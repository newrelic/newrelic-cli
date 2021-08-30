## newrelic profile add

Add a new profile

### Synopsis

Add a new profile

The add command creates a new profile for use with the New Relic CLI.
API key and region are required. An Insights insert key is optional, but required
for posting custom events with the `newrelic events`command.


```
newrelic profile add [flags]
```

### Examples

```
newrelic profile add --profile <profile> --region <region> --apiKey <apiKey> --insightsInsertKey <insightsInsertKey> --accountId <accountId> --licenseKey <licenseKey>
```

### Options

```
  -y, --acceptDefaults             suppress prompts and accept default values
      --apiKey string              your personal API key
  -h, --help                       help for add
      --insightsInsertKey string   your Insights insert key
      --licenseKey string          your license key
  -r, --region string              the US or EU region
```

### Options inherited from parent commands

```
  -a, --accountId int    the account ID to use. Can be overridden by setting NEW_RELIC_ACCOUNT_ID
      --debug            debug level logging
      --format string    output text format [JSON, Text, YAML] (default "JSON")
      --plain            output compact text
      --profile string   the authentication profile to use
      --trace            trace level logging
```

### SEE ALSO

* [newrelic profile](newrelic_profile.md)	 - Manage the authentication profiles for this tool

