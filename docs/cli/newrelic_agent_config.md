## newrelic agent config

Configuration utilities/helpers for New Relic agents

### Synopsis

Configuration utilities/helpers for New Relic agents

### Examples

```
newrelic agent config obfuscate --value <config_value> --key <obfuscation_key>
```

### Options

```
  -h, --help   help for config
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

* [newrelic agent](newrelic_agent.md)	 - Utilities for New Relic Agents
* [newrelic agent config migrateV3toV4](newrelic_agent_config_migrateV3toV4.md)	 - migrate V3 configuration to V4 configuration format
* [newrelic agent config obfuscate](newrelic_agent_config_obfuscate.md)	 - Obfuscate a configuration value using a key

