## newrelic nerdstorage document write

Write a NerdStorage document.

### Synopsis

Write a NerdStorage document

Write a NerdStorage document.  Valid scopes are ACCOUNT, ENTITY, and USER.
ACCOUNT scope requires a valid account ID and ENTITY scope requires a valid entity
GUID.  A valid Nerdpack package ID is required.


```
newrelic nerdstorage document write [flags]
```

### Examples

```

  # Account scope
  newrelic nerdstorage document write --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol --documentId myDoc --document '{"field": "myValue"}'

  # Entity scope
  newrelic nerdstorage document write --scope ENTITY --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --entityId MjUyMDUyOHxFUE18QVBQTElDQVRJT058MjE1MDM3Nzk1 --collection myCol --documentId myDoc --document '{"field": "myValue"}'

  # User scope
  newrelic nerdstorage document write --scope USER --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --collection myCol --documentId myDoc --document '{"field": "myValue"}'

```

### Options

```
  -c, --collection string   the collection name to write the document to
  -o, --document string     the document to be written, in JSON format (default "{}")
  -d, --documentId string   the document ID
  -e, --entityGuid string   the entity GUID
  -h, --help                help for write
  -p, --packageId string    the external package ID
  -s, --scope string        the scope to write the document to (default "USER")
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

* [newrelic nerdstorage document](newrelic_nerdstorage_document.md)	 - Read, write, and delete NerdStorage documents.

