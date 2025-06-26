package template

const (
	ServerAddCommandsWithStartCmdOptions = "server.AddCommandsWithStartCmdOptions"

	RollkitV0XStartHandler = "rollserv.StartHandler"
	RollkitV1XStartHandler = "abciserver.StartHandler"
	rollkitV1MigrateCmd    = "abciserver.MigrateToRollkitCmd()"
)

const (
	GoExecPackage = "github.com/rollkit/go-execution-abci"
	GoExecVersion = "v0.2.1-0.20250625133753-4c5a41d10330" // TODO(@julienrbrt): use tag when available

	RollkitPackage = "github.com/rollkit/rollkit"
	RollkitVersion = "v0.14.2-0.20250625130707-6ad581a97a16" // TODO(@julienrbrt): use tag when available
)
