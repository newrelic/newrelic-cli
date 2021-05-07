## newrelic utils terraform

Tools for working with Terraform

### Synopsis

Tools for working with Terraform

The terraform commands can be used for generating Terraform HCL for simple observability
as code use cases.


### Examples

```
cat terraform.json | newrelic utils terraform dashboard --label my_dashboard_resource
```

### Options

```
  -h, --help   help for terraform
```

### Options inherited from parent commands

```
      --format string   output text format [JSON, Text, YAML] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic utils](newrelic_utils.md)	 - Various utility methods
* [newrelic utils terraform dashboard](newrelic_utils_terraform_dashboard.md)	 - Generate HCL for the newrelic_one_dashboard resource

