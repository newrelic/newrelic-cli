## newrelic nerdstorage document get

Retrieve a NerdStorage document.

### Synopsis

Retrieve a NerdStorage document

Retrieve a NerdStorage document.  Valid scopes are ACCOUNT, ENTITY, and USER.
ACCOUNT scope requires a valid account ID and ENTITY scope requires a valid entity
GUID.  A valid Nerdpack package ID is required.


```
newrelic nerdstorage document get [flags]
```

### Examples

```

  # Account scope
  newrelic nerdstorage document get --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol --documentId myDoc

  # Entity scope
  newrelic nerdstorage document get --scope ENTITY --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --entityId MjUyMDUyOHxFUE18QVBQTElDQVRJT058MjE1MDM3Nzk1  --collection myCol --documentId myDoc

  # User scope
  newrelic nerdstorage document get --scope USER --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --collection myCol --documentId myDoc

```

### Options

```
  -a, --accountId int       the account ID
  -c, --collection string   the collection name to get the document from
  -d, --documentId string   the document ID
  -e, --entityGuid string   the entity GUID
  -h, --help                help for get
  -p, --packageId string    the external package ID
  -s, --scope string        the scope to get the document from (default "USER")
```

### Options inherited from parent commands

```
      --format string   output text format [YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic nerdstorage document](newrelic_nerdstorage_document.md)	 - Read, write, and delete NerdStorage documents.

