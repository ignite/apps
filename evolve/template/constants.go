package template

const (
	ServerAddCommandsWithStartCmdOptions = "server.AddCommandsWithStartCmdOptions"

	RollkitV0XStartHandler = "rollserv.StartHandler"
	EvolveV1XStartHandler  = "abciserver.StartHandler"
	evolveV1MigrateCmd     = "abciserver.MigrateToEvolveCmd()"
)

const (
	GoExecPackage = "github.com/evstack/ev-abci"
	GoExecVersion = "v0.3.1-0.20250818124535-74c4793fcbee" // https://github.com/evstack/ev-abci/pull/202

	EvNodePackage = "github.com/evstack/ev-node"
	EvNodeVersion = "v1.0.0-beta.2.0.20250818104031-a7dbf779dbfe" // main until tag

	EvNodeDaCmd     = "github.com/evstack/ev-node/da/cmd/local-da"
	EvNodeDaVersion = "v1.0.0-beta.1"
)
