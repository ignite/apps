package template

const (
	ServerAddCommandsWithStartCmdOptions = "server.AddCommandsWithStartCmdOptions"

	RollkitV0XStartHandler = "rollserv.StartHandler"
	EvolveV1XStartHandler  = "abciserver.StartHandler"
	evolveV1MigrateCmd     = "abciserver.MigrateToEvolveCmd()"
)

const (
	EvABCIPackage = "github.com/evstack/ev-abci"
	EvABCIVersion = "v0.4.0"

	EvNodePackage = "github.com/evstack/ev-node"
	EvNodeVersion = "v1.0.0-beta.3"

	GoDataStorePackageFork = "github.com/celestiaorg/go-datastore"
	GoDataStoreVersionFork = "v0.0.0-20250801131506-48a63ae531e4"
	GoDataStorePackage     = "github.com/ipfs/go-datastore"

	EvNodeDaCmd = "github.com/evstack/ev-node/da/cmd/local-da"
)
