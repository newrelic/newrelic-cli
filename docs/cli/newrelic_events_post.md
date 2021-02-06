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
  -a, --accountId int   the account ID to create the custom event in
  -e, --event string    a JSON-formatted event payload to post (default "{}")
  -h, --help            help for post
```

### Options inherited from parent commands

```
      --format string   output text format [JSON, Text, YAML] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic events](newrelic_events.md)	 - Send custom events to New Relic

