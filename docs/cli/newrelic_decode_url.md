## newrelic decode url

Decodes NR1 URL Strings 

```
newrelic decode url [flags]
```

### Examples

```
newrelic decode url -p="pane" -s="entityId" https://one.newrelic.com/launcher/nr1-core.home?pane=eyJuZXJkbGV0SWQiOiJucjEtY29yZS5ob21lLXNjcmVlbiJ9&platform[accountId]=1
```

### Options

```
  -h, --help            help for url
  -p, --param string    the query parameter you want to decode
  -s, --search string   the search key you want returned
```

### Options inherited from parent commands

```
  -a, --accountId int    trace level logging
      --debug            debug level logging
      --format string    output text format [JSON, Text, YAML] (default "JSON")
      --plain            output compact text
      --profile string   the authentication profile to use
      --trace            trace level logging
```

### SEE ALSO

* [newrelic decode](newrelic_decode.md)	 - Decodes NR1 URL Strings 

