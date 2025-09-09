package template

const (
	ServerAddCommandsWithStartCmdOptions = "server.AddCommandsWithStartCmdOptions"

	RollkitV0XStartHandler = "rollserv.StartHandler"
	EvolveV1XStartHandler  = "abciserver.StartHandler"
	evolveV1MigrateCmd     = "abciserver.MigrateToEvolveCmd()"
)

const (
	EvABCIPackage = "github.com/evstack/ev-abci"
	EvABCIVersion = "v0.3.1-0.20250818181501-f014411689fd" // main until tag

	EvNodePackage = "github.com/evstack/ev-node"
	EvNodeVersion = "v1.0.0-beta.2.0.20250818133040-d096a24e7052" // main until tag

	EvNodeDaCmd     = "github.com/evstack/ev-node/da/cmd/local-da"
	EvNodeDaVersion = "v1.0.0-beta.1"
)
