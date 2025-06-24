package template

const (
	ServerAddCommandsWithStartCmdOptions = "server.AddCommandsWithStartCmdOptions"

	RollkitV0XStartHandler = "rollserv.StartHandler"
	RollkitV1XStartHandler = "abciserver.StartHandler"
	rollkitV1MigrateCmd    = "abciserver.MigrateToRollkitCmd()"
)

const (
	GoExecPackage = "github.com/rollkit/go-execution-abci"
	GoExecVersion = "v0.2.1-0.20250624195837-8bbd403344b4" // TODO(@julienrbrt): use tag when available

	RollkitPackage = "github.com/rollkit/rollkit"
	RollkitVersion = "v0.14.2-0.20250624162917-8bb1ba71af39" // TODO(@julienrbrt): use tag when available
)
