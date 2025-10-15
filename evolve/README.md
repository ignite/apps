# Evolve

This Ignite App is aimed to extend [Ignite CLI](https://github.com/ignite/cli) and bootstrap the development of a [Evolve](https://rollkit.dev) network.

## Prerequisites

- Ignite CLI version v28.9.0 or greater.
- Knowledge of blockchain development (Cosmos SDK).

## Usage

```sh
ignite s chain gm --address-prefix gm --minimal --no-module
cd gm
ignite app install -g github.com/ignite/apps/evolve@latest
ignite evolve add
ignite chain build --skip-proto
ignite evolve init # only for genesis chains. Otherwise follow the migration steps.
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

### Migrations

If you want to migrate your running chain to Evolve, first scaffold the migrations commands on your CometBFT chain:

```sh
ignite evolve add-migrate
```

This will add the migration module to your chain. Then add manually a chain migration in the upgrade handler to add this new module and submit a gov proposal to initiate the validator set migration.

Once the chain has halted, run the migration command on each node:

```sh
gmd evolve-migrate
```

You are ready to integrate Evolve! Follow the [1](#Usage) steps to add it to your chain.

Learn more about Evolve and Ignite in their respective documentation:

- <https://docs.ignite.com>
- <https://ev.xyz/>
