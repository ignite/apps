
  

# Ignite App: # Chain Info Example

This example demonstrates how to create an Ignite App to gather information about a chain and build it using the Ignite CLI.
  

## Installation

  

### Within Project Directory

  

To use the `chain-info` app within your project, execute the following command inside the project directory:

  

```bash

ignite  app  install  github.com/ignite/apps/examples/chain-info

```

  

The app will be available only when running `ignite` inside the project directory.

  

### Globally

  

To use the `chain-info` app globally, execute the following command:

  

```bash

ignite  app  install  -g  github.com/ignite/apps/examples/chain-info

```

  

This command will compile the app and make it immediately available to the `ignite` command lists.

  

## Requirements

  

- Go (version 1.16 or higher)

- Ignite CLI (version 28.1.1 or higher)

  

## How it Works

The `chain-info` app allows users to gather information about a chain in the current directory and build it using the Ignite CLI.

-   `main.go`: Defines the main functionality of the `chain-info` Ignite App. It registers the app with the Ignite CLI, defines the app's metadata, and handles the execution of commands.
    
-   `cmd/cmd.go`: Contains the definition of the main command and its subcommands. It specifies the available commands and their descriptions.
    
-   `cmd/info.go`: Contains the execution logic for the `info` subcommand. It gathers information about the chain, including its version, app path, configuration path, initialization status, and binary file.
    
-   `cmd/build.go`: Contains the execution logic for the `build` subcommand. It builds the chain app in the current directory using Ignite helper functions.
- 
## Integration Test

The integration test for the `chain-info` Ignite App ensures that the app functions correctly with the Ignite CLI. The test performs the following steps:

1.  Installs the `chain-info` app locally.
2.  Runs the `info` subcommand to gather information about the chain in the current directory.
3.  Runs the `build` subcommand to build the chain app in the current directory.
4.  Verifies the output of the `info` subcommand to ensure it displays the correct information about the chain.
5.  Verifies the output of the `build` subcommand to ensure it indicates successful chain building.

To run the integration test, execute the following command:

```bash
go test -v ./integration
```
