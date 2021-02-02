## NewRelic decode

Decode a New Relic String to reveal information about the page application

### Synopsis

Use the decode command to print out information encrypted within the URL. 

```
newrelic decode [flags]
```

### Examples

```
newrelic decode https://one.newrelic.com/launcher/nr1-core.home?pane=eyJuZXJkbGV0SWQiOiJucjEtY29yZS5ob21lLXNjcmVlbiJ9&platform[accountId]=1

```

### Options

```
  -h, --help   help for Decoded
```

### Options inherited from parent commands

```
      --format string   output text format [Text, YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic](newrelic.md)	 - The New Relic CLI

