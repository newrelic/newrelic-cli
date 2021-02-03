## NewRelic decode

Decode a New Relic String to reveal information about the page application


### Synopsis

Use the decode url command to print out information encrypted within the URL. 

```
newrelic decode url [flags]
```

### Examples

```
newrelic decode url -p="pane" -j="entityId" https://one.newrelic.com/launcher/nr1-core.home?pane=eyJuZXJkbGV0SWQiOiJucjEtY29yZS5ob21lLXNjcmVlbiJ9&platform[accountId]=1

```

### Options

```
  -h, --help    Help for Decoded
  -p, --param   The param you want to decode
  -j, --json    The json element you want returned
```

### Options inherited from parent commands

```
      --format string   output text format [Text, YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [New Relic](newrelic.md)	 - The New Relic CLI

