## newrelic agent config obfuscate

Obfuscate a configuration value using a key

### Synopsis

Obfuscate a configuration value using a key.  The obfuscated value
should be placed in the Agent configuration or in an environment variable." 


```
newrelic agent config obfuscate [flags]
```

### Examples

```
newrelic agent config obfuscate --value <config_value> --key <obfuscation_key>
```

### Options

```
  -h, --help           help for obfuscate
  -k, --key string     the key to use when obfuscating the clear-text value
  -v, --value string   the value, in clear text, to be obfuscated
```

### Options inherited from parent commands

```
      --format string   output text format [YAML, JSON, Text] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic agent config](newrelic_agent_config.md)	 - Configuration utilities/helpers for New Relic agents

