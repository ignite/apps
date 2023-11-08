# Ignite Apps

Official [Ignite CLI](https://ignite.com/cli) apps repository.

Each directory in the root of this repo must be a Go module containing
an Ignite App package, each one with its own `go.mod` file.

## Developer instruction

- Clone this repo locally.
- Scaffold your app: `ignite app scaffold my-app`
- Add the folder to the `go.work`.
- Add your cobra commands into `debug/main.go` and the module replace to the `debug/go.mod` for a easy debug.
- Add the plugin: `ignite app add -g ($GOPATH)/src/github.com/ignite/apps/my-app`
- Test with Ignite.
