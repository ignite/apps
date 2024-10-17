# Fee Abstraction

The Fee Abstraction app is an extension for the [Ignite CLI](https://github.com/ignite/cli), designed to help developers seamlessly integrate the [Fee Abstraction](https://github.com/osmosis-labs/fee-abstraction) module from Osmosis Labs into their blockchain projects.

This app extends the `ignite scaffold chain` command by adding a `--fee-abstraction` flag, which automatically incorporates the Fee Abstraction module into your chain.

## Features

- Effortless Integration: Easily integrate the Fee Abstraction module into your blockchain with a single command.

## Prerequisites

- Ignite CLI version v28.5.2 or greater.
  - Or migrate your chain using the [Ignite migration guides](https://docs.ignite.com/migration).
- Knowledge of blockchain development and the Fee Abstraction module.

## Installation

To install the Fee Abstraction app, it need to be globally, run the following command:

```shell
ignite app install -g github.com/ignite/apps/fee-abstraction
```

## Usage

- Now you can scaffold your chain using the fee abstraction module:

```shell
ignite s chain mars --fee-abstraction
```

This command will scaffold a new chain named `mars` with the Fee Abstraction module already integrated.
