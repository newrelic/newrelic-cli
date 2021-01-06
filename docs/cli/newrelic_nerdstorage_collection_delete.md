## newrelic nerdstorage collection delete

Delete a NerdStorage collection.

### Synopsis

Delete a NerdStorage collection

Delete a NerdStorage collection.  Valid scopes are ACCOUNT, ENTITY, and USER.
ACCOUNT scope requires a valid account ID and ENTITY scope requires a valid entity
GUID.  A valid Nerdpack package ID is required.


```
newrelic nerdstorage collection delete [flags]
```

### Examples

```

  # Account scope
  newrelic nerdstorage collection delete --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol

  # Entity scope
  newrelic nerdstorage collection delete --scope ENTITY --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --entityId MjUyMDUyOHxFUE18QVBQTElDQVRJT058MjE1MDM3Nzk1  --collection myCol

  # User scope
  newrelic nerdstorage collection delete --scope USER --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --collection myCol

```

### Options

```
  -a, --accountId int       the account ID
  -c, --collection string   the collection name to delete the document from
  -e, --entityGuid string   the entity GUID
  -h, --help                help for delete
      --packageId string    the external package ID (default "p")
  -s, --scope string        the scope to delete the document from (default "USER")
```

### Options inherited from parent commands

```
      --format string   output text format [YAML, JSON, Text] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic nerdstorage collection](newrelic_nerdstorage_collection.md)	 - Read, write, and delete NerdStorage collections.

