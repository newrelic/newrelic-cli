## newrelic events post

Post a custom event to New Relic

### Synopsis

Post a custom event to New Relic

The post command accepts an account ID and JSON-formatted payload representing a
custom event to be posted to New Relic. These events once posted can be queried
using NRQL via the CLI or New Relic One UI.
The accepted payload requires the use of an `eventType`field that
represents the custom event's type.


```
newrelic events post [flags]
```

### Examples

```
newrelic events post --accountId 12345 --event '{ "eventType": "Payment", "amount": 123.45 }'
```

### Options

```
  -e, --event string   a JSON-formatted event payload to post (default "{}")
  -h, --help           help for post
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

* [newrelic events](newrelic_events.md)	 - Send custom events to New Relic

