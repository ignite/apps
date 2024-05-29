# RollKit

This Ignite App is aimed to extend [Ignite CLI](https://github.com/ignite/cli) and bootstrap the development of a [RollKit](https://rollkit.dev) rollup.

## Prerequisites

* Ignite CLI version v28.4.0 or greater.
* Knowledge of blockchain development (Cosmos SDK).

## Usage

```sh
ignite s chain gm --address-prefix gm
cd gm
ignite app install -g github.com/ignite/apps/rollkit@rollkit/v0.1.0
ignite rollkit add
ignite chain build
ignite rollkit init
```

Learn more about Rollkit and Ignite in their respective documentation:

* <https://docs.ignite.com>
* <https://rollkit.dev/>
