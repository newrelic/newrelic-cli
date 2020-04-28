# Extensions

The New Relic CLI is an extensible platform for building on the New Relic Platform.


## Use Cases

Basic use cases and examples for the main types of interactions a plugin might want to do.

1. **Do something locally**

   Create a new nerdpack:

   ```
   # Existing command in nr1 CLI
   nr1 create

   # Example command as Extension
   newrelic one create
   ```

1. **Do something against the New Relic APIs**

   Make a NRQL query:

   ```
   # Existing command in nr1 CLI
   nr1 nrql --account $NEW_RELIC_ACCOUNT_ID --query 'FROM Transaction SELECT count(*)'

   # Example command as Extension
   newrelic nrql
   ```

1. **Do something against a third-party API**

   Run a NRQL query and send the resulting link to a Slack channel

   ```
   #Run NRQL and post to slack channel
   query($account:Int!, $query:Nrql!) {
     actor {
       account(id: $account) {
         nrql(query: $query) {
           staticChartUrl(chartType: LINE)
         }
       }
     }
   }
   {
     "account": $NEW_RELIC_ACCOUNT_ID,
     "query":"FROM Transaction SELECT count(*) TIMESERIES SINCE 1 hour ago"
   }

   Post "data"."actor"."account"."nrql"."staticChartUrl" link to Slack

   newrelic slack-me-a-chart \
     --query "FROM Transaction SELECT count(*) TIMESERIES SINCE 1 hour ago"
   ```

1. **Intercept output and reformat it**

   We can't implement all output formats, so why not allow an extension to process the data and format it?

   ```
   # CSV would have to be a plugin name, pre-installed and registered
   newrelic nerdgraph query --query "" --variables "" --output csv

   accountId,name
   0,badwolf
   ```

## Architecture

### CLI

#### Responsibilities

* Manage the installed extensions
  * List available extensions
  * Download and Install
  * Execute
  * Remove
* Enable communication to known NR APIs
* Handle basic I/O with the user (terminal)
* Stateful config required by the extensions


#### Exported Handlers

* Input (stdin)
* Output (stdout)
* Logging (stderr)
* Configuration
  * Load
  * Save
* HTTP
  * Abstract handling of
     * Retries
     * timeouts
     * proxies
  * New Relic API(s)
     * Uses profile
     * Region
     * Authentication headers
  * Arbitrary Endpoint
     * Raw request, unmodified



```
+--+ CLI
|  |
|  +--+ handlers
|  |  |
|  |  +-- Register your handler endpoint
|  |  |
|  |  |
|  |  +-- New Relic API Handler (HTTP Requests to NR)
|  |  |
|  |  +--
|
|
```

### Extension

Example extension schema definition file.  For configs, I like YAML vs JSON, since you can have comments in YAML.

```
+--+ extension.yaml file
   |
   +-- protocol_version
   |
   +--+ metadata
   |  |
   |  +-- name
   |  +-- version
   |  +-- description
   |  +-- checksum / validation of some sort (terraform does this how?)
   |  +--+ requires
   |     +-- cli_version
   |     +-- extensions []
   |     +-- runtime
   |
   +--+ handlers
   |  |
   |  +--+ input
   |  |  +-- name
   |  |
   |  +--+ output
   |     +-- name
   |
   +--+ commands { }
      |
      +--+ "mycommand"
         +-- name
         +-- help / description
         +-- aliases []
         +--+ flags { }
         |  +-- name
         |  +-- help / description
         |  +-- data type
         |
         +-- commands { }
```
