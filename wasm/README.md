# **Wasm**

Wasm is an Ignite App developed to extend the [Ignite CLI](https://github.com/ignite/cli), enabling developers to integrate CosmWasm smart contracts into their blockchain projects with ease. This app adds commands for adding and configuring CosmWasm support within your chain.

The app adds `ignite wasm` commands to add a [CosmWasm](https://cosmwasm.com/) integration into a chain.

## **Features**

- Easily add CosmWasm integration to your blockchain.
- Streamline the process of configuring Wasm for your chain.
- Develop and test CosmWasm smart contracts in your Ignite CLI project.

## **Prerequisites**

- Ignite CLI version v28.3.0 or greater. 
    - Or migrate your chain using the [Ignite migration guides](https://docs.ignite.com/migration).
- Knowledge of blockchain development and smart contracts.

## **Installation**

- Install the Wasm app:

```shell
ignite app install -g github.com/ignite/apps/wasm
```

- You must scaffold a new chain with version `v28.2.1` or greater. 

- Navigate to your chain's directory and execute the following command to add Wasm support:

```shell
ignite wasm add
```

This command integrates Wasm into your chain's code and configuration. If your chain configuration does not exist yet (for non-initiated chains), you'll need to add the Wasm configuration manually:

```shell
ignite wasm config
```

Remember, all commands should be executed within your chain directory.

## **Developer instruction**

To contribute to the Wasm app or use a local version, follow these steps:

```shell
git clone github.com/ignite/apps && cd apps/wasm
```

Add the app to the global config:

```shell
ignite app add -g /absolute/path/to/app/wasm # or use $(pwd)
```

The `ignite wasm` command is now available with the local version of the app.

1. **Develop and Test:**
- Make changes to the app code as needed.
- Run **`ignite wasm`** to recompile the app and test your changes.

## **Support and Contributions**

For support and contributions, please visit the [Wasm App GitHub repository](https://github.com/ignite/apps/wasm). 
We welcome contributions from the community, including bug reports, feature requests, and code contributions.
