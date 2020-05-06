# New Relic CLI Extensions

## Overview
The New Relic CLI provides an extension framework as a way to integrate new commands into its core command set.  Extensions can be built in any language.

## Getting started

This section will guide you through creating an extension for the New Relic CLI.  At its core, an extension is a repository that contains an `extension.yml` in its root, and a way to invoke a process that will start the extension.

### Execution modes
There are two execution models for CLI extensions, API mode and binary execution mode.

#### API mode
In this execution model, the extension is responsible for communicating with the CLI host via an HTTP API.  The API allows for rich interaction with the user (e.g. interactive prompts), conduits for displaying log messages, warnings, errors, and command output, and a means of communicating to the New Relic backend via Nerdgraph queries.

The interaction begins with a call to [`/command`](#command), which will return the command that the user has invoked along with any provided flags.  The extension can then carry out its work and pass its output to the user via the [`/out`](#out) endpoint.  See the [CLI API reference](#cli-api) section below for more details on communicating with the CLI API.

#### Binary execution mode
In this execution model, the extension simply receives its command and arguments as command line arguments during invocation.  IO from the extension is redirected to the user's terminal.

When using `binary` mode, the extension's entrypoint (see [entrypoint](#entrypoint) above) will be invoked with additional arguments representing the command to be executed and a collection of flag value pairs, in the format `<ENTRYPOINT> <COMMAND> <FLAG_1> <VALUE_1> <FLAG_2> <VALUE_2> <FLAG_N> <VALUE_N>`.  An example invocation:

```
node index.js hello name Shelly
```

## Extension manifest

Extensions need to expose some information about the functionality that they provide in the form of an `extension.yml` file.  The manifest is used by the CLI host to invoke the extension and generate help screens for the user.

### Field reference

* [`name`](#name)
* [`entrypoint`](#entrypoint)
* [`description`](#description)
* [`execution_type`](#execution_type)
* [`commands`](#commands)
  * [`name`](#commands.name)
  * [`short_description`](#commands.short_description)
  * [`long_description`](#commands.long_description)
  * [`example`](#commands.example)
  * [`parent`](#commands.parent)
  * [`flags`](#commands.flags)
    * [`name`](#commands.flags.name)
    * [`type`](#commands.flags.type)
    * [`default`](#commands.flags.default)
    * [`usage`](#commands.flags.usage)
    * [`shorthand`](#commands.flags.shorthand)

#### `name`
* **Type**: `string`
* **Required**: true

The name of the extension.  This will map to a subcommand in the CLI that the user can call alongside the core subcommands.

#### `version`
* **Type**: `string`
* **Required**: true

The version of this extension, in `MAJOR.MINOR.PATCH` format.

#### `entrypoint`
* **Type**: `string`
* **Required**: true

A command that can be invoked by the CLI host to start the extension.  Entrypoint commands are executed from within the root directory of the extension.

#### `description`
* **Type**: `string`
* **Required**: true

A description of the extension as it will appear in the CLI help. This will be presented in the CLI help along with the top-level extension command.

#### `execution_type`
* **Type**: `string`
* **Required**: false
* **Default**: `api`

The execution model this extension will use. Options are `api` and `binary` (See [Execution Modes](#plugin-modes) for details on execution models).

#### `commands`
* **Type**: `array of hashes`
* **Required**: true

The commands this extension provides.

#### `commands.name`
* **Type**: `string`
* **Required**: true

The name of the command.  This will map to a child command under the extension's root command that can be invoked by the user.

#### `commands.short_description`
* **Type**: `string`
* **Required**: false

The short description of the command.  This will be presented in the CLI help along with the command.

#### `commands.long_description`
* **Type**: `string`
* **Required**: false

The long description of the command.  This will be presented in the CLI help along with the command.

#### `commands.example`
* **Type**: `string`
* **Required**: false

A string that contains a sample invocation of the command.  This will be presented in the CLI help along with the command.

#### `commands.parent`
* **Type**: `array of strings`
* **Required**: false

*Advanced use only*.  This argument allows the extension author to add child commands to other commands in the CLI command tree besides the top-level extension command defined above. Core commands cannot be overwritten via this method.

#### `commands.flags`
* **Type**: `array of hashes`
* **Required**: false
* **Default**: `[]`

The flags this command provides to the user.

#### `commands.flags.name`
* **Type**: `string`
* **Required**: true

The name of the flag.  The user can supply this as a long option to the CLI command being invoked.

#### `commands.flags.type`
* **Type**: `string`
* **Required**: false
* **Default**: `string`

The type of the flag's value.  Options are `string`, `bool`, or `int`.

#### `commands.flags.default`
* **Type**: `string`, `bool`, or `int` (See [commands.flags.type](#commands.flags.type) above.)
* **Required**: false

The default value for this flag.

#### `commands.flags.usage`
* **Type**: `string`
* **Required**: false

The usage description for this flag.  This will be presented in the CLI help along with the flag.

#### `commands.flags.shorthand`
* **Type**: `string`
* **Required**: false

A one-character shorthand for the flag name. The user can supply this as a short option to the CLI command being invoked.

#### Example
```
name: hello-world
short: hello world commands
entrypoint: node index.js
commands:
  - name: hello
    short_description: say hello
    long_description: This command allows the user to generate a hello world message.
    example: |
      newrelic hello -n Shelly

      > Hello Shelly!
    flags:
      - name: name
        shorthand: n
        type: string
        default: friend
        usage: your name
```

## CLI API
If the extension is making use of the `api` mode (see [Execution modes](#plugin-modes) above), it can communicate with the CLI host via an HTTP API.

### Authentication
Upon invocation, the extension will be provided with a URI for connecting to the host CLI as well as a session token for authentication purposes.  The authentication token should be provided in the `X-CLI-Auth` header for all requests back to the CLI host.

When using `api` mode, the extension's entrypoint (see [entrypoint](#entrypoint) above) will be invoked with two additional arguments representing the host HTTP server's URI and a session token, in the format `<ENTRYPOINT> <URI> <TOKEN>`.  An example invocation:

```
node index.js localhost:23855 65f6b964-a19f-4698-b2db-4e171420439b
```

Requests that do not provide the auth token will result in 401 Unauthorized responses.

### Endpoint reference

* [`Command`](#command)
* [`Out`](#out)
* [`Log`](#log)
* [`Prompt`](#prompt)
* [`Query`](#query)
* [`Fail`](#fail)


#### `Command`

Used to collect user context for the current command execution.  Contains the command to be run as well as any flag values supplied by the user.  This is the first call the extension should make after initializing.

**URL** : `/command`

**Method** : `GET`

**Response**

```json
{
    "command": "hello",
    "flags": {
      "name": "Shelly"
    }
}
```

#### `Out`

Used to send output to the user's terminal.

**URL** : `/out`

**Method** : `POST`

**Request body**
```json
{
  "value": "Hello Shelly!"
}
```


#### `Log`

Used to send an informational log message to the user's terminal.  Supported levels are `info`, `warn`, and `error`.

**URL** : `/log`

**Method** : `POST`

**Request body**
```json
{
  "level": "info",
  "message": "Generating hello message..."
}
```

#### `Prompt`

Used to collect user input via an interactive prompt.  `options`, `default`, and `placeholder` can be provided but are not required.

**URL** : `/prompt`

**Method** : `POST`

**Request body**
```json
{
  "prompt": "Who would you like to say hello to?",
  "options": ["Shelly", "Tim", "Ramses"],
  "default": "Shelly",
  "placeholder": "Shelly",
}
```

**Response**

```json
{
  "value": "Shelly"
}
```

#### `Query`

Used to execute an authenticated request against the NerdGraph API.  The host CLI will use the current user's CLI profile to determine how to authenticate to the NerdGraph backend.

**URL** : `/query`

**Method** : `POST`

**Request body**
```json
{
  "query": "query($accountId: Int!){ actor { account(id: $accountId) { name } } } ",
  "variables": {
    "accountId": 12345
  }
}
```

**Response**

```json
{
  "data": {
    "actor": {
      "account": {
        "name": "My Organization's New Relic Account"
      }
    }
  }
}
```

#### `Fail`

Used to indicate to the host CLI that the extension has failed.

**URL** : `/fail`

**Method** : `POST`

## Publishing and installing extensions
To allow the CLI host to install the extension, it must be publicly available via GitHub.  It must also include an `extension.yml` in the project root as defined above.

Users can then install a plugin via the `newrelic plugins add` command:

```
> newrelic plugins add github.com/newrelic/hello-world-extension
```

Commands can be removed via the `newrelic plugins remove` command based on the name defined in their manifest:

```
> newrelic plugins remove hello-world
```

To view a list of installed plugins, the user can use the `newrelic plugins list` command:

```
> newrelic plugins list

+-------------+---------+-------------------------------------------+
|    Name     | Version |                Repository                 |
+-------------+---------+-------------------------------------------+
| hello-world |   0.0.1 | github.com/newrelic/hello-world-extension |
+-------------+---------+-------------------------------------------+

```

The user can use the `newrelic plugins upgrade` command to upgrade an installed plugin:

```
> newrelic plugins upgrade
```