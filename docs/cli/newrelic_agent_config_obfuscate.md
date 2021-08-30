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
  -a, --accountId int    the account ID to use. Can be overridden by setting NEW_RELIC_ACCOUNT_ID
      --debug            debug level logging
      --format string    output text format [JSON, Text, YAML] (default "JSON")
      --plain            output compact text
      --profile string   the authentication profile to use
      --trace            trace level logging
```

### SEE ALSO

* [newrelic agent config](newrelic_agent_config.md)	 - Configuration utilities/helpers for New Relic agents

