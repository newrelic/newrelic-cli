## newrelic agent config migrateV3toV4

migrate V3 configuration to V4 configuration format

### Synopsis

migrate V3 configuration to V4 configuration format

```
newrelic agent config migrateV3toV4 [flags]
```

### Examples

```
newrelic integrations config migrateV3toV4 --pathDefinition /file/path --pathConfiguration /file/path --pathOutput /file/path
```

### Options

```
  -h, --help                       help for migrateV3toV4
      --overwrite                  if set ti true and pathOutput file exists already the old file is removed 
  -c, --pathConfiguration string   path configuration
  -d, --pathDefinition string      path definition
  -o, --pathOutput string          path output
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

