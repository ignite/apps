# debug

This package is a helper that you could use to debug your Ignite App by running it as a standalone binary independently of Ignite CLI.

## How to use

```shell
go run tools/debug/main.go <COMMAND>
```

e.g:

```shell
go run tools/debug/main.go explorer <ARGS>
```

```shell
go run tools/debug/main.go hermes <ARGS>
```

## Developer instruction

- Replace the app repo for a local folder into the `tools/debug/go.mod`.

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

### Caveat

The app doesn't support debugging interactions with the CLI. This method allows debugging running App commands independently of the CLI, which means that Ignite doesn't know the PID of the App, so it won't be able to attach to it, and because of that, debugged Apps won't be able to communicate with Ignite to for example use the Client API.
We will soon support dynamic debugging in Ignite CLI.
