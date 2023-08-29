# relayer

`relayer` is a plugin developed for [Ignite CLI](https://github.com/ignite/cli).

The plugin adds `ignite relayer` commands that allow launching new Cosmos blockchains by interacting with the Ignite Chain to coordinate with validators.

The plugin is integrated into Ignite CLI by default.

## Developer instruction

- clone this repo locally
- Run `ignite plugin add -g /absolute/path/to/plugin/relayer` to add the plugin to global config
- `ignite relayer` command is now available with the local version of the plugin.

Then repeat the following loop:

- Hack on the plugin code
- Rerun `ignite relayer` to recompile the plugin and test
