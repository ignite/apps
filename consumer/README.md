# Consumer app

Consumer app an Ignite CLI app for ICS consumer chain.

It can be used to run 2 specific tasks, identified by 2 arguments passed to `ExecutedCommand.Args` field:
- `writeGenesis`: write the consumer genesis 
- `isInitialized`: verify that a consumer chain is properly initialized

The only goal of using an app for these tasks is to avoid the interchain-security dependency inside Ignite CLI.
