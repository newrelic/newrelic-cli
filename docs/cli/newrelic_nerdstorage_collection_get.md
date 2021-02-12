## newrelic nerdstorage collection get

Retrieve a NerdStorage collection.

### Synopsis

Retrieve a NerdStorage collection

Retrieve a NerdStorage collection.  Valid scopes are ACCOUNT, ENTITY, and USER.
ACCOUNT scope requires a valid account ID and ENTITY scope requires a valid entity
GUID.  A valid Nerdpack package ID is required.


```
newrelic nerdstorage collection get [flags]
```

### Examples

```

  # Account scope
  newrelic nerdstorage collection get --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol

  # Entity scope
  newrelic nerdstorage collection get --scope ENTITY --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --entityId MjUyMDUyOHxFUE18QVBQTElDQVRJT058MjE1MDM3Nzk1  --collection myCol

  # User scope
  newrelic nerdstorage collection get --scope USER --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --collection myCol

```

### Options

```
  -a, --accountId int       the account ID
  -c, --collection string   the collection name to get the document from
  -e, --entityGuid string   the entity GUID
  -h, --help                help for get
  -p, --packageId string    the external package ID
  -s, --scope string        the scope to get the document from (default "USER")
```

### Options inherited from parent commands

```
      --format string   output text format [YAML, JSON, Text] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic nerdstorage collection](newrelic_nerdstorage_collection.md)	 - Read, write, and delete NerdStorage collections.

