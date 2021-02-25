## newrelic edge trace-observer

Interact with New Relic Edge trace observers.

### Synopsis

Interact with New Relic Edge trace observers
	
	A trace observer is a configuration that enables infinite tracing for an account.
	Once enabled, infinite tracing observes 100% of your application traces, then
	provides visualization for the most actionable data so you can investigate and
	solve issues faster.

### Examples

```
newrelic edge trace-observer list --accountId <accountID>
```

### Options

```
  -a, --accountId int   A New Relic account ID
  -h, --help            help for trace-observer
```

### Options inherited from parent commands

```
      --format string   output text format [YAML, JSON, Text] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic edge](newrelic_edge.md)	 - Interact with New Relic Edge
* [newrelic edge trace-observer create](newrelic_edge_trace-observer_create.md)	 - Create a New Relic Edge trace observer.
* [newrelic edge trace-observer delete](newrelic_edge_trace-observer_delete.md)	 - Delete a New Relic Edge trace observer.
* [newrelic edge trace-observer list](newrelic_edge_trace-observer_list.md)	 - List the New Relic Edge trace observers for an account.

