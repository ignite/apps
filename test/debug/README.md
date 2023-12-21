# debug

This package is a helper for you can debug your app outside the Ignite using devel or another debuger.

## How to use

```shell
go run test/debug/main.go <COMMAND>
```

e.g:
```shell
go run test/debug/main.go explorer <ARGS>
```

```shell
go run test/debug/main.go hermes <ARGS>
```

## Developer instruction

- Replace the app repo for a local folder into the `test/debug/go.mod`.
```go.mod
replace (
	github.com/ignite/apps/explorer => ../../explorer
	github.com/ignite/apps/hermes => ../../hermes
	github.com/ignite/apps/<MY-APP> => ../../<MY-APP> 
)
```

- Add the command to be debugged.
```go
rootCmd.AddCommand(
    explorer.NewExplorer(),
    hermes.NewHermes(),
    myapp.NewCommand() // <--- Add the new command
)
```