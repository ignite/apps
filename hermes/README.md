# relayer

`relayer` is a app developed for [Ignite CLI](https://github.com/ignite/cli).

The app adds `ignite relayer hermes` commands to allow IBC communication between two different blockchain apps using [Hermes relayer](https://github.com/informalsystems/hermes).

## How to use

- Run both chain to be relayed.

- add the hermes relayer app from remote:
```shell
ignite app install -g github.com/ignite/apps/hermes
```

- or clone the repo and add the hermes relayer app local:
```shell
ignite app install -g $GOPATH/src/github.com/ignite/apps/hermes
```

- configure the relayer:
```shell
ignite relayer hermes configure [chain-a-id] [chain-a-rpc] [chain-a-grpc] [chain-b-id] [chain-b-rpc] [chain-b-grpc] [flags]
```
e.g.:
```shell
ignite relayer hermes configure "mars-1" "http://localhost:26649" "http://localhost:9082" "venus-1" "http://localhost:26659" "http://localhost:9092"
```

- start the relayer
```shell
ignite relayer hermes start [chain-a-id] [chain-b-id] [flags]
```
e.g.:
```shell
ignite relayer hermes start "mars-1" "venus-1"
```


## Developer instruction

- clone this repo locally
- Run `ignite app add -g /absolute/path/to/app/relayer` to add the app to global config
- `ignite relayer` command is now available with the local version of the app.

Then repeat the following loop:

- Hack on the app code
- Rerun `ignite relayer` to recompile the app and test
