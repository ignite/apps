# Wasm

`wasm` is an app developed for [Ignite CLI](https://github.com/ignite/cli).

The app adds `ignite wasm` commands to add a [CosmWasm](https://cosmwasm.com/) integration into a chain.

## How to use

- Install the Wasm app:
```shell
ignite app install -g github.com/ignite/apps/wasm
```

- You must scaffold a new chain with version `v28.2.1` or greater. Or migrate your chain using the [Ignite migration guides](https://docs.ignite.com/migration).

- Now you can add a Wasm support to your chain by running this command into the chain directory:
```shell
ignite wasm add
```

- The command will automatically add the Wasm integration into your code and the Wasm config into your chain config. But if your chain is not initiated yet, the chain config does not exist, so you need to add the Wasm config later running:
```shell
ignite wasm add
```

_All commands should be run in the chain directory._


## Developer instruction

- clone this repo locally
- Run `ignite app add -g /absolute/path/to/app/wasm` to add the app to global config
- `ignite wasm` command is now available with the local version of the app.

Then repeat the following loop:

- Hack on the app code
- Rerun `ignite wasm` to recompile the app and test
