## newrelic nrql droprules

Retrieves NRQL drop rules

### Synopsis

Retrieve NRQL drop rules

The droprules command will fetch a list of NRQL drop rules linked to your account.

```
newrelic nrql droprules [flags]
```

### Examples

```
newrelic nrql droprules
newrelic nrql droprules --limit 2
```

### Options

```
  -h, --help        help for droprules
  -l, --limit int   number of droprules to return (default 10)
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

* [newrelic nrql](newrelic_nrql.md)	 - Commands for interacting with the New Relic Database

