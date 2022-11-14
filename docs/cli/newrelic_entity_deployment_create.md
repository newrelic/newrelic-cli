## newrelic entity deployment create

Create a deployment marker for a given entity

### Synopsis

Create a deployment marker for a given entity

*NOTE:* This feature is in Limited Preview and not yet available to all customers.

The create command marks a deployment for the given New Relic entity.


```
newrelic entity deployment create [flags]
```

### Examples

```
newrelic entity deployment create --guid <GUID> --version <0.0.1> --changelog 'what changed' --commit '12345e' --deepLink <link back to deployer> --deploymentType 'BASIC' --description 'about' --timestamp <1668446197100> --user 'jenkins-bot'
```

### Options

```
  -h, --help                    help for create
  -g, --guid string             the entity GUID for the deployment. guid is required.
  -v, --version string          the version of the deployed software, for example, something like v1.1. version is required.
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
