# Wasm

Wasm is an Ignite App developed to extend the [Ignite CLI](https://github.com/ignite/cli), enabling developers to integrate CosmWasm smart contracts into their blockchain projects with ease. This app adds commands for adding and configuring CosmWasm support within your chain.

The app adds `ignite wasm` commands to add a [CosmWasm](https://cosmwasm.com/) integration into a chain.

## Features

- Easily add CosmWasm integration to your blockchain.
- Streamline the process of configuring Wasm for your chain.
- Develop and test CosmWasm smart contracts in your Ignite CLI project.

## Prerequisites

- Ignite CLI version v28.3.0 or greater.
  - Or migrate your chain using the [Ignite migration guides](https://docs.ignite.com/migration).
- Knowledge of blockchain development and smart contracts.

## Installation

```shell
ignite app install -g github.com/ignite/apps/wasm
```

> [!IMPORTANT]  
> Wasmd uses CGO. Which means library for compiling C must be installed.
> While on most linux distribution or macos those are insatalled by default,
> users of WSL or Ubuntu must run the following: `sudo apt-get install build-essential`

## Usage

- Navigate to your chain's directory and execute the following command to add Wasm support:

```shell
ignite wasm add
```

This command integrates Wasm into your chain's code and configuration. If your chain configuration does not exist yet (for non-initiated chains), you'll need to add the Wasm configuration manually:

```shell
ignite chain init
ignite wasm config
```

Remember, all commands should be executed within your chain directory.

## Configuration

In order to configure CosmWasm as permissioned in your chain, you can add the following configuration to your chain's `config.yaml` file:

```yaml
genesis:
  app_state:
    wasm:
      params:
        code_upload_access:
          addresses: []
          permission: "Nobody"
        instantiate_default_permission: "Nobody"
```

Read [CosmosWasm docs](https://github.com/CosmWasm/wasmd/blob/21b048d54e395ff9168e5c3037356a73797500ba/x/wasm/Governance.md?plain=1#L27-L48) for more information on the `code_upload_access` configuration.
