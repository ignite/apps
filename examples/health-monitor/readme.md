
  

# Ignite App: Health Monitor Example

The `health-monitor` Ignite App is an example demonstrating how to implement a health monitoring application with Ignite.
  

## Installation

  

### Within Project Directory

  

To use the `health-monitor` app within your project, execute the following command inside the project directory:

  

```bash

ignite  app  install  github.com/ignite/apps/examples/health-monitor

```

  

The app will be available only when running `ignite` inside the project directory.

  

### Globally

  

To use the `health-monitor` app globally, execute the following command:

  

```bash

ignite  app  install  -g  github.com/ignite/apps/examples/health-monitor

```

  

This command will compile the app and make it immediately available to the `ignite` command lists.

  

## Requirements

  

- Go (version 1.21 or higher)

- Ignite CLI (version 28.1.1 or higher)

  

## How it Works
The `health-monitor` Ignite App implements a health monitoring application that retrieves and displays information about the status of a running chain. It consists of several files:

-   `main.go`: Contains the main code for the app, including the implementation of the `Manifest` method, which defines the app's name and commands, and the `Execute` method, which handles the execution of the `monitor` subcommand based on the provided flags and arguments.
    
-   `cmd/cmd.go`: Contains the definition of the main command and its subcommands, including their names, descriptions, and flags.
    
-   `cmd/monitor.go`: Contains the execution logic for the `monitor` subcommand, which retrieves and prints the health status of a running chain.
- 
## Integration Test
The integration test for the `health-monitor` Ignite App ensures that the app works correctly with the Ignite CLI. The test performs the following steps:

1.  Installs the `health-monitor` app locally.
2.  Runs the `monitor` subcommand with custom flags to specify the RPC address and refresh duration.
3.  Verifies the output of the `monitor` subcommand.

To run the integration test, execute the following command:

```bash
go test -v ./integration
```
