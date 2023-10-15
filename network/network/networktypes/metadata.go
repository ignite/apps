package networktypes

// Version current version of the network plugin
const Version = "1"

// Cli holds information about the CLI used to interact with the chain
type Cli struct {
	Version string `json:"version"`
}

// Metadata is an object that contains the metadata of a chain
// the metadata represents generic data set by the coordinator
// these information can be formatted and interpreted by the CLI for specific purposes
type Metadata struct {
	Cli Cli `json:"cli"`
}

// IsCurrentVersion checks if the version of the CLI is equal to the current version
func (m Metadata) IsCurrentVersion() bool {
	return m.Cli.Version == Version
}
