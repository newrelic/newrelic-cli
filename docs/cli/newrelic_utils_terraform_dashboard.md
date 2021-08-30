## newrelic utils terraform dashboard

Generate HCL for the newrelic_one_dashboard resource

### Synopsis

Generate HCL for the newrelic_one_dashboard resource

This command generates HCL configuration for newrelic_one_dashboard resources from
exported JSON documents.  For more detail on exporting dashboards to JSON, see
https://docs.newrelic.com/docs/query-your-data/explore-query-data/dashboards/manage-your-dashboard/#dash-json

Input can be sourced from STDIN per the provided example, or from a file using the --file option.
Output will be sent to STDOUT by default but can be redirected to a file with the --out option.


```
newrelic utils terraform dashboard [flags]
```

### Examples

```
cat terraform.json | newrelic utils terraform dashboard --label my_dashboard_resource
```

### Options

```
  -f, --file string      a file that contains exported dashboard JSON
  -h, --help             help for dashboard
  -l, --label string     the resource label to use when generating resource HCL
  -o, --out string       the file to send the generated HCL to
  -w, --shiftWidth int   the indentation shift with of the output (default 2)
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

* [newrelic utils terraform](newrelic_utils_terraform.md)	 - Tools for working with Terraform

