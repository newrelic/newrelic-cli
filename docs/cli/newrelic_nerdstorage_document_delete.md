## newrelic nerdstorage document delete

Delete a NerdStorage document.

### Synopsis

Delete a NerdStorage document

Delete a NerdStorage document.  Valid scopes are ACCOUNT, ENTITY, and USER.
ACCOUNT scope requires a valid account ID and ENTITY scope requires a valid entity
GUID.  A valid Nerdpack package ID is required.


```
newrelic nerdstorage document delete [flags]
```

### Examples

```

  # Account scope
  newrelic nerdstorage document delete --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol --documentId myDoc

  # Entity scope
  newrelic nerdstorage document delete --scope ENTITY --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --entityId MjUyMDUyOHxFUE18QVBQTElDQVRJT058MjE1MDM3Nzk1 --collection myCol --documentId myDoc

  # User scope
  newrelic nerdstorage document delete --scope USER --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --collection myCol --documentId myDoc

```

### Options

```
  -c, --collection string   the collection name to delete the document from
  -d, --documentId string   the document ID
  -e, --entityGuid string   the entity GUID
  -h, --help                help for delete
  -p, --packageId string    the external package ID
  -s, --scope string        the scope to delete the document from (default "USER")
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

* [newrelic nerdstorage document](newrelic_nerdstorage_document.md)	 - Read, write, and delete NerdStorage documents.

