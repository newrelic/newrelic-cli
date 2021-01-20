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
newrelic profile add --name <profileName> --region <region> --apiKey <apiKey> --insightsInsertKey <insightsInsertKey> --accountId <accountId> --licenseKey <licenseKey>
```

### Options

```
      --accountId int              your account ID
      --apiKey string              your personal API key
  -h, --help                       help for add
      --insightsInsertKey string   your Insights insert key
      --licenseKey string          your license key
  -n, --name string                unique profile name to add
  -r, --region string              the US or EU region
```

### Options inherited from parent commands

```
      --format string   output text format [Text, YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic profile](newrelic_profile.md)	 - Manage the authentication profiles for this tool

