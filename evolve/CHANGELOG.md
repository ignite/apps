# Evolve App Changelog

## Unreleased

## [`v0.5.0`](https://github.com/ignite/apps/releases/tag/evolve/v0.5.0)

- []() Remove `--start` and `--migrate` flags from `evolve add` command. Split into `evolve add` and `evolve add-migrate` commands.
- [#236](https://github.com/ignite/apps/pull/236) Add `--start` flag to `evolve add` command to optionally disable addition of the start command.

## [`v0.4.3`](https://github.com/ignite/apps/releases/tag/evolve/v0.4.3)

- [#233](https://github.com/ignite/apps/pull/233) Bump dependencies.

## [`v0.4.2`](https://github.com/ignite/apps/releases/tag/evolve/v0.4.2)

- [#232](https://github.com/ignite/apps/pull/232) Fix node syncing by bumping dependencies.

## [`v0.4.1`](https://github.com/ignite/apps/releases/tag/evolve/v0.4.1)

- [#229](https://github.com/ignite/apps/pull/229) Update dependencies.

## [`v0.4.0`](https://github.com/ignite/apps/releases/tag/evolve/v0.4.0)

- [#220](https://github.com/ignite/apps/pull/220) Rename app, flags and commands from `rollkit` to `evolve` following the rebranding of Rollkit to Evolve.
- [#223](https://github.com/ignite/apps/pull/223) Wire rollback command.
- [#225](https://github.com/ignite/apps/pull/225) Init namespace with chain id on `init`.

## [`v0.3.0`](https://github.com/ignite/apps/releases/tag/rollkit/v0.3.0)

- [#112](https://github.com/ignite/apps/pull/112) Use default command instead cobra commands.
- [#192](https://github.com/ignite/apps/pull/192) Upgrade Rollkit to `v1.x` and use [`go-execution-abci`](https://github.com/evstack/ev-abci) instead of [`cosmos-sdk-starter`](https://github.com/rollkit/cosmos-sdk-starter).
  - Note, if you already have rollkit installed, you need to redo the wiring manually.
- [#194](https://github.com/ignite/apps/pull/194) Uniform the flag addition with other apps.
- [#209](https://github.com/ignite/apps/pull/209) Wire `MigrateToRollkitCmd` in the scaffolded app.
- [#212](https://github.com/ignite/apps/pull/212) Add `--migrate` flag to `rollkit add` command to scaffold modules and migration helpers for CometBFT to Rollkit migration.
- [#212](https://github.com/ignite/apps/pull/212) Add `edit` command to edit an existing genesis file without overwriting it (compared to `init`).

## [`v0.2.3`](https://github.com/ignite/apps/releases/tag/rollkit/v0.2.3)

- Dependency bumps

## [`v0.2.2`](https://github.com/ignite/apps/releases/tag/rollkit/v0.2.2)

- [#168](https://github.com/ignite/apps/pull/168) Update expected template.

## [`v0.2.1`](https://github.com/ignite/apps/releases/tag/rollkit/v0.2.1)

- [#106](https://github.com/ignite/apps/pull/106) Improve bonded tokens setting.

## [`v0.2.0`](https://github.com/ignite/apps/releases/tag/rollkit/v0.2.0)

- [#92](https://github.com/ignite/apps/pull/92) Add `rollkit init` command.

## [`v0.1.0`](https://github.com/ignite/apps/releases/tag/rollkit/v0.1.0)

- First release of the Rollkit app compatible with Ignite >= v28.x.y
