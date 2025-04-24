# RollKit

This Ignite App is aimed to extend [Ignite CLI](https://github.com/ignite/cli) and bootstrap the development of a [RollKit](https://rollkit.dev) rollup.

## Prerequisites

* Ignite CLI version v28.9.0 or greater.
* Knowledge of blockchain development (Cosmos SDK).

## Usage

```sh
ignite s chain gm --address-prefix gm --minimal
cd gm
ignite app install -g github.com/ignite/apps/rollkit@latest
ignite rollkit add
ignite chain build
ignite rollkit init
```

Then start `local-da` or use Celestia mainnet as data availability layer.

```sh
# go install github.com/rollkit/rollkit/da/cmd/local-da@latest
git clone github.com/rollkit/rollkit --depth 1
cd rollkit/da/cmd/local-da
go run .
```

Finally, run the rollup node:

```sh
gmd start --rollkit.node.aggregator
```

Learn more about Rollkit and Ignite in their respective documentation:

* <https://docs.ignite.com>
* <https://rollkit.dev/>
