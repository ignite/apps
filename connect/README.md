# Connect

This Ignite App extends [Ignite CLI](https://github.com/ignite/cli) to let a user interact with any Cosmos SDK based chain.

## Installation

```shell
ignite app install -g github.com/ignite/apps/connect
```

### Usage

* Discover available chains

```shell
ignite connect discover
```

* Add a chain to interact with

```shell
ignite connect add atomone
```

* (Or) Add a local chain to interact with

```shell
ignite connect add simapp localhost:9090
```

* List all connected chains

```shell
ignite connect
```

* Remove a connected chain

```shell
ignite connect rm atomone
```
