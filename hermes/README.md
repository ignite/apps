# relayer

`relayer` is a app developed for [Ignite CLI](https://github.com/ignite/cli).

The app adds `ignite relayer` commands that allow launching new Cosmos blockchains by interacting with the Ignite Chain to coordinate with validators.

The app is integrated into Ignite CLI by default.

## Developer instruction

- clone this repo locally
- Run `ignite app add -g /absolute/path/to/app/relayer` to add the app to global config
- `ignite relayer` command is now available with the local version of the app.

Then repeat the following loop:

- Hack on the app code
- Rerun `ignite relayer` to recompile the app and test
