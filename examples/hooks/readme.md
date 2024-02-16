
  

# Ignite App: Hooks Example

This example demonstrates how to use hooks in an Ignite App to extend the functionality of the Ignite CLI.
  

## Installation

  

### Within Project Directory

  

To use the `hooks` app within your project, execute the following command inside the project directory:

  

```bash

ignite  app  install  github.com/ignite/apps/examples/hooks

```

  

The app will be available only when running `ignite` inside the project directory.

  

### Globally

  

To use the `hooks` app globally, execute the following command:

  

```bash

ignite  app  install  -g  github.com/ignite/apps/examples/hooks

```

  

This command will compile the app and make it immediately available to the `ignite` command lists.

  

## Requirements

  

- Go (version 1.21 or higher)

- Ignite CLI (version 28.1.1 or higher)

  

## How it Works

The `hooks` app registers hooks on two Ignite CLI commands: `ignite chain build` and `ignite chain serve`. These hooks are triggered before and after the execution of the respective commands and provide additional information about the chain running at the RPC address in the current directory where Ignite is running.

The `main.go` file defines the main functionality of the `hooks` Ignite App. Here's a breakdown of its components:

-   `Manifest`: This method defines the app's metadata, including its name (`hooks`) and the list of hooks it provides. Two hooks are registered: `chain-build` and `chain-serve`, attached to specific Ignite CLI commands (`ignite chain build` and `ignite chain serve` respectively). Additionally, this method retrieves the list of commands defined in the `cmd/cmd.go` file using `cmd.GetCommands()`.
    
-   `Execute`: Handles the execution of the app. It prints a message indicating how to use the app by running either `ignite chain build` or `ignite chain serve`.
    
-   `ExecuteHookPre`: Called before the execution of a hook. It prints a message indicating that the hook is about to be executed and retrieves information about the chain (if any) running at the RPC address in the current directory.
    
-   `ExecuteHookPost`: Called after the successful execution of a hook. It prints a message indicating that the hook has been executed successfully.
    
-   `ExecuteHookCleanUp`: Called after the execution of a hook, regardless of the result. It performs any necessary cleanup tasks related to the hook execution.
    
-   `main`: Initializes and serves the Ignite App using the HashiCorp Go plugin framework. It configures the handshake and plugins, where the `hooks` plugin is registered with the app instance created from the `app` struct.
    

Overall, `main.go` sets up the core functionality of the `hooks` Ignite App, defining hooks, handling their execution, and serving the app using the HashiCorp Go plugin framework.
