# Ignite Airdrop App

Official [Ignite CLI](https://ignite.com/cli) apps repository.

Each directory in the root of this repo must be a Go module containing
an Ignite App package, each one with its own `go.mod` file.

### How to test

Run the raw airdrop command to generate the raw data
```shell
go run cmd/debug/debug.go airdrop raw testdata/genesis.json 2> raw-snapshot.json
```

Run the process command to generate the claim records
```shell
go run cmd/debug/debug.go airdrop process testdata/config.yml raw-snapshot.json 2> snapshot.json
```

Generate the new genesis json with the claimable state based on the output genesis
```shell
go run cmd/debug/debug.go airdrop genesis testdata/config.yml raw-snapshot.json testdata/genesis.json
```


Or you can only use the generate command to run all four above with one command
```shell
go run cmd/debug/debug.go airdrop generate testdata/config.yml testdata/genesis.json testdata/genesis.json
```

