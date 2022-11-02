## newrelic entity deployment create

Create a deployment marker for a given entity

### Synopsis

Create a deployment marker for a given entity

The create command marks a deployment for the given New Relic entity.


```
newrelic entity deployment create [flags]
```

### Examples

```
newrelic entity deployment create --guid <GUID> --version <1.0.0>
```

### Options

```
  -h, --help                    help for create
  -g, --guid string             the entity GUID to create change tracker
  -v, --version string          the tag names to add to the entity
      --changelog string        a URL for the changelog or list of changes if not linkable
      --commit string           the commit identifier, for example, a Git commit SHA
      --deepLink string         a link back to the system generating the deployment
      --deploymentType string   type of deployment, one of BASIC, BLUE_GREEN, CANARY, OTHER, ROLLING or SHADOW
      --description string      a description of the deployment
      --groupID string          string that can be used to correlate two or more events
  -t  --timestamp int64         the start time of the deployment, the number of milliseconds since the Unix epoch, defaults to now       
  -u  --user string             username of the deployer or bot
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

* [newrelic entity deployment](newrelic_entity_deployment.md) - Track deployments for a New Relic entity 

