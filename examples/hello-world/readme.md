
# Ignite App: Hello World Example

The "hello-world" Ignite App is a simple example demonstrating how to create a basic Ignite App that prints "Hello, world!" when executed. 

## Installation

### Within Project Directory

To use the "hello-world" app within your project, execute the following command inside the project directory:

```bash
ignite app install github.com/ignite/apps/examples/hello-world
```

The app will be available only when running `ignite` inside the project directory.

### Globally

To use the "hello-world" app globally, execute the following command:

```bash
ignite app install -g github.com/ignite/apps/examples/hello-world
```

This command will compile the app and make it immediately available to the `ignite` command lists.

## Requirements

- Go (version 1.16 or higher)
- Ignite CLI (version 28 or higher)

## How it Works

The "hello-world" Ignite App is a simple example of an Ignite App that implements a basic command to print "Hello, world!". The app consists of two main files:

- `main.go`: Contains the main code for the app, including the implementation of the `Manifest` method, which defines the app's name and commands, and the `Execute` method, which defines the execution logic for the command.

- `cmd/cmd.go`: Contains the definition of the command, including its name and description.

When the app is executed using the Ignite CLI, it prints "Hello, world!" to the console.

## Integration Test

The integration test for the "hello-world" Ignite App ensures that the app works correctly with the Ignite CLI. The test performs the following steps:

1. Installs the "hello-world" app locally.
2. Runs the app.
3. Asserts that the output is "Hello, world!".

To run the integration test, execute the following command:

```bash
go test -v ./integration_test
```
