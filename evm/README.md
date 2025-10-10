# EVM

This Ignite App is aimed to extend [Ignite CLI](https://github.com/ignite/cli) and bootstrap the development of an [EVM](https://github.com/cosmos/evm)-enabled Cosmos SDK blockchain.

## Prerequisites

- Ignite CLI version v29.5.0 or greater (no minimal chain).
- Knowledge of blockchain development (Cosmos SDK).

## Usage

```sh
ignite s chain gm --address-prefix gm --no-module
cd gm
ignite app install -g github.com/ignite/apps/evm
ignite evm add
```
