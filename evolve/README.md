# Evolve

This Ignite App is aimed to extend [Ignite CLI](https://github.com/ignite/cli) and bootstrap the development of a [Evolve](https://rollkit.dev) network.

## Prerequisites

* Ignite CLI version v28.9.0 or greater.
* Knowledge of blockchain development (Cosmos SDK).

## Usage

```sh
ignite s chain gm --address-prefix gm --minimal --no-module
cd gm
ignite app install -g github.com/ignite/apps/evolve@latest
ignite evolve add
ignite chain build --skip-proto
ignite evolve init
```

Then start `local-da` or use Celestia mainnet as data availability layer.

```sh
cd gm
go tool github.com/evstack/ev-node/da/cmd/local-da
```

Finally, run the network:

```sh
gmd start --rollkit.node.aggregator
```

Learn more about Evolve and Ignite in their respective documentation:

* <https://docs.ignite.com>
* <https://rollkit.dev/>
