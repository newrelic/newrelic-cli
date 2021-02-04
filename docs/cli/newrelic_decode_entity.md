## NewRelic decode

Decode a New Relic String to reveal information about the page application


### Synopsis

Use the decode url command to print out information encrypted within the URL. 

```
newrelic decode entity [flags]
```

### Examples

```
newrelic decode entity -k=ID  MXxBUE18QVBQTElDQVRJT058Mzk4NDkyNDQw 
```

### Options

```
  -h, --help           Help for Decoded
  -k, --key string     The key you want returned from an entity

```

### Options inherited from parent commands

```
      --format string   output text format [Text, YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [NewRelic](newrelic.md)	 - The New Relic CLI


