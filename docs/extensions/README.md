# Extensions

The New Relic CLI is an extensible platform for building on the New Relic Platform.

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
