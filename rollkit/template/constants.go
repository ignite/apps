package template

const (
	ServerAddCommandsWithStartCmdOptions = "server.AddCommandsWithStartCmdOptions"

	RollkitV0XStartHandler = "rollserv.StartHandler"
	RollkitV1XStartHandler = "abciserver.StartHandler"
	rollkitV1MigrateCmd    = "abciserver.MigrateToRollkitCmd()"
)

const (
	GoExecPackage = "github.com/rollkit/go-execution-abci"
	GoExecVersion = "v0.3.0"

	RollkitPackage = "github.com/rollkit/rollkit"
	RollkitVersion = "v1.0.0-beta.2"

	RollkitDaCmd     = "github.com/rollkit/rollkit/da/cmd/local-da"
	RollkitDaVersion = "v1.0.0-beta.1"
)
