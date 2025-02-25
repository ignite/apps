# Ignite App: Spinner Example

The "spinner" Ignite App is a simple example demonstrating how to create a basic Ignite App that prints the spinner
interaction.

## Installation

### Within Project Directory

To use the "spinner" app within your project, execute the following command inside the project directory:

```bash
ignite app install github.com/ignite/apps/examples/spinner
```

The app will be available only when running `ignite` inside the project directory.

### Globally

To use the "spinner" app globally, execute the following command:

```bash
ignite app install -g github.com/ignite/apps/examples/spinner
```

This command will compile the app and make it immediately available to the `ignite` command lists.

## Requirements

- Go (version 1.16 or higher)
- Ignite CLI (version 28.1.1 or higher)

## How it Works

The "spinner" Ignite App is a simple example of an Ignite App that implements a basic command to print the spinner
interaction. The app consists of two main files:

- `main.go`: Contains the main code for the app, including the implementation of the `Manifest` method, which defines
  the app's name and commands, and the `Execute` method, which defines the execution logic for the command.

- `cmd/cmd.go`: Contains the definition of the command, including its name and description.

When the app is executed using the Ignite CLI, it prints the spinner interaction to the console.
