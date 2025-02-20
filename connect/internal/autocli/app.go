package autocli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ignite/apps/connect/chains"
	"github.com/ignite/apps/connect/internal/autocli/flag"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
)

// Builder manages options for building CLI commands.
type Builder struct {
	// flag.Builder embeds the flag builder and its options.
	flag.Builder

	// Config is the config of the chain from the connect app
	Config *chains.ChainConfig

	// GetClientConn specifies how CLI commands will resolve a grpc.ClientConnInterface
	// from a given context.
	GetClientConn func(*cobra.Command) (grpc.ClientConnInterface, error)

	// AddQueryConnFlags adds flags to query commands
	AddQueryConnFlags func(*cobra.Command)

	// AddTxConnFlags adds flags to transaction commands
	AddTxConnFlags func(*cobra.Command)

	// Cdc is the codec to use for encoding and decoding messages.
	Cdc codec.Codec
}

// EnhanceRootCommand enhances the root command with the provided module options.
//
// ModuleOptions are autocli options to be used for modules. They are gotten from
// the reflection service.
func EnhanceRootCommand(
	rootCmd *cobra.Command,
	builder *Builder,
	moduleOptions map[string]*autocliv1.ModuleOptions,
) error {
	if err := builder.Validate(); err != nil {
		return err
	}

	queryCmd, err := builder.BuildQueryCommand(rootCmd.Context(), moduleOptions)
	if err != nil {
		return err
	}

	msgCmd, err := builder.BuildMsgCommand(rootCmd.Context(), moduleOptions)
	if err != nil {
		return err
	}

	rootCmd.AddCommand(queryCmd, msgCmd)

	return nil
}
