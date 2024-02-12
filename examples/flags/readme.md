
  

# Ignite App: Flags Example

The `flags` Ignite App is a simple example demonstrating how to use command-line flags and arguments in Ignite applications.
  

## Installation

  

### Within Project Directory

  

To use the `flags` app within your project, execute the following command inside the project directory:

  

```bash

ignite  app  install  github.com/ignite/apps/examples/flags

```

  

The app will be available only when running `ignite` inside the project directory.

  

### Globally

  

To use the `flags` app globally, execute the following command:

  

```bash

ignite  app  install  -g  github.com/ignite/apps/examples/flags

```

  

This command will compile the app and make it immediately available to the `ignite` command lists.

  

## Requirements

  

- Go (version 1.21 or higher)

- Ignite CLI (version 28.1.1 or higher)

  

## How it Works

  

The `flags` Ignite App demonstrates the use of command-line flags and arguments in Ignite applications. It consists of several files:
-   `main.go`: Contains the main code for the app, including the implementation of the `Manifest` method, which defines the app's name and commands, and the `Execute` method, which handles the execution of subcommands based on the provided arguments and flags.
    
-   `cmd/cmd.go`: Contains the definition of the main command and its subcommands, including their names, descriptions, and flags.
    
-   `cmd/hello.go` and `cmd/cowsay.go`: Contains the execution logic for the `hello` and `cowsay` subcommands, respectively. These files demonstrate how to access and use flags passed to the subcommands.
  

The `hello` subcommand simply prints a greeting message to the console, including the name specified by the `--name` flag. For example:
 The `cowsay` subcommand uses the `Neo-cowsay` library to generate a ASCII art of a cow saying hello, including the name specified by the `--name` flag.
 

## Integration Test

  

The integration test for the `flags` Ignite App ensures that the app works correctly with the Ignite CLI. The test performs the following steps:

  

1.  Installs the `flags` app locally.

2.  Runs the `hello` subcommand with a custom `name` flag.

3.  Runs the `cowsay` subcommand with a custom `name` flag.

5.  Verifies the output of the `hello` and `cowsay` subcommands.

  

To run the integration test, execute the following command:

  

```bash

go test -v ./integration

```
