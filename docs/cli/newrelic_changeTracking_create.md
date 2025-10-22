## newrelic changeTracking create

Create a New Relic change tracking event

### Synopsis

Create a New Relic change tracking event

This command allows you to create a change tracking event for a New Relic entity, supporting all fields in the Change Tracking GraphQL API schema for the changeTrackingCreateEvent mutation. For more information on each field, visit: https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#change-tracking-event-mutation

```
newrelic changeTracking create [flags]
```

### Examples

```
newrelic changeTracking create \
  --entitySearch "name = 'MyService' AND type = 'SERVICE'" \
  --category Deployment \
  --type Basic \
  --description "Deployed version 1.2.3 to production" \
  --version "1.2.3" \
  --changelog "https://github.com/myorg/myservice/releases/tag/v1.2.3" \
  --commit "abc123def456" \
  --user "ci-cd-bot"
```

### Options

Required fields:
```
  --entitySearch        Entity search query (e.g., name = 'MyService' AND type = 'SERVICE'). See our docs on 'entitySearch.query' under 'Required attributes' for more detailed examples: https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#required-fields
  --category            Category of event (e.g. Deployment, Feature Flag, Operational, etc.)
  --type                Type of event (e.g. Basic, Rollback, Server Reboot, etc.)
```

For Deployment events, the following are required/supported:
```
  --version             Version of the deployment (required)
  --changelog           Changelog for the deployment (URL or text)
  --commit              Commit hash for the deployment
  --deepLink            Deep link URL for the deployment
```

For Feature Flag events, the following are required/supported:
```
  --featureFlagId       ID of the feature flag (required)
```

Other supported fields:
```
  --description         Description of the event
  --user                Username of the actor or bot
  --groupId             String to correlate two or more events
  --shortDescription    Short description for the event
  --customAttributes    Custom attributes: use '-' for STDIN, '{...}' for inline JS object, or provide a file path
  --validationFlags     Comma-separated list of validation flags (e.g. ALLOW_CUSTOM_CATEGORY_OR_TYPE, FAIL_ON_FIELD_LENGTH, FAIL_ON_REST_API_FAILURES)
  --timestamp           Time of the event (milliseconds since Unix epoch, defaults to now). Can not be more than 24 hours in the past or future
```

Custom attributes can be provided in three ways:
  1. From STDIN by passing '-' (e.g. `echo  '{cloud_vendor: "vendor_name", region: "us-east-1", isProd: true, instances: 2}' | newrelic changeTracking create ... --customAttributes -`)
  2. As an inline JS object starting with '{' (e.g. `--customAttributes '{cloud_vendor: "vendor_name", region: "us-east-1", isProd: true, instances: 2}`)
  3. As a file path (e.g. `--customAttributes ./attrs.js`)
The JS object format must use unquoted keys and values of type string, boolean, or number. Example: `{cloud_vendor: "vendor_name", region: "us-east-1", isProd: true, instances: 2}`
Validation is performed before sending to the API. Keys must be valid JS identifiers, and values must be string, boolean, or number.

For more information, see: https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#change-tracking-event-mutation

### Options inherited from parent commands

```
  -a, --accountId int       the account ID to use. Can be overridden by setting NEW_RELIC_ACCOUNT_ID
      --debug               debug level logging
      --format string       output text format [JSON, Text, YAML] (default "JSON")
      --plain               output compact text
      --profile string      the authentication profile to use
      --trace               trace level logging
```

### SEE ALSO

* [newrelic](newrelic.md)	 - The New Relic CLI
