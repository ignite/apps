package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cosmossdk.io/core/address"
	"github.com/cosmos/cosmos-sdk/client"
	sdkflags "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/client/v2/autocli/flag"
	"github.com/ignite/apps/connect/chains"
	"github.com/ignite/apps/connect/internal"
)

func AppHandler(ctx context.Context, name string, cfg *chains.ChainConfig, args ...string) (*cobra.Command, error) {
	chainCmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Commands for %s chain", name),
	}

	if len(args) > 0 {
		chainCmd.SetArgs(args)
	}

	conn, err := chains.NewConn(name, cfg)
	if err != nil {
		return nil, err
	}

	if err := conn.Load(ctx); err != nil {
		return nil, err
	}

	addressCodec, validatorAddressCodec, consensusAddressCodec := setupAddressPrefixesAndCodecs(cfg.Bech32Prefix)

	builder := &autocli.Builder{
		Builder: flag.Builder{
			TypeResolver:          &dynamicTypeResolver{conn},
			FileResolver:          conn.ProtoFiles,
			AddressCodec:          addressCodec,
			ValidatorAddressCodec: validatorAddressCodec,
			ConsensusAddressCodec: consensusAddressCodec,
		},
		GetClientConn: func(cmd *cobra.Command) (grpc.ClientConnInterface, error) {
			return conn.Connect()
		},
		AddQueryConnFlags: func(cmd *cobra.Command) {
			sdkflags.AddQueryFlagsToCmd(cmd)
			sdkflags.AddKeyringFlags(cmd.Flags())
		},
		AddTxConnFlags: sdkflags.AddTxFlagsToCmd,
	}

	// add comet commands
	cometCmds := cmtservice.NewCometBFTCommands()
	conn.ModuleOptions[cometCmds.Name()] = cometCmds.AutoCLIOptions()

	// add autocli commands
	appOpts := &autocli.AppOptions{
		ModuleOptions:         conn.ModuleOptions,
		AddressCodec:          addressCodec,
		ValidatorAddressCodec: validatorAddressCodec,
		ConsensusAddressCodec: consensusAddressCodec,
	}

	// keyring config
	k, err := internal.NewKeyring(chainCmd.Flags(), addressCodec, cfg.Bech32Prefix)
	if err != nil {
		return nil, err
	}

	// create client context
	clientCtx := client.Context{}.
		WithKeyring(k).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithViper("")

	// add to root command (autocli expects it there)
	chainCmd.SetContext(context.WithValue(
		context.Background(),
		client.ClientContextKey,
		&clientCtx,
	))

	err = appOpts.EnhanceRootCommandWithBuilder(chainCmd, builder)
	if err != nil {
		return nil, err
	}

	return chainCmd, nil
}

// setupAddressPrefixesAndCodecs returns the address codecs for the given bech32 prefix.
// Additionally it sets the address prefix for the sdk.Config.
func setupAddressPrefixesAndCodecs(prefix string) (
	address.Codec,
	address.Codec,
	address.Codec,
) {
	// set address prefix for sdk.Config
	var (
		// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key.
		bech32PrefixAccPub = prefix + sdk.PrefixPublic
		// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address.
		bech32PrefixValAddr = prefix + sdk.PrefixValidator + sdk.PrefixOperator
		// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key.
		bech32PrefixValPub = bech32PrefixValAddr + sdk.PrefixPublic
		// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address.
		bech32PrefixConsAddr = prefix + sdk.PrefixValidator + sdk.PrefixConsensus
		// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key.
		bech32PrefixConsPub = bech32PrefixConsAddr + sdk.PrefixPublic
	)

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(prefix, bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(bech32PrefixValAddr, bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(bech32PrefixConsAddr, bech32PrefixConsPub)
	config.Seal()

	return addresscodec.NewBech32Codec(prefix),
		addresscodec.NewBech32Codec(bech32PrefixValAddr),
		addresscodec.NewBech32Codec(bech32PrefixConsAddr)
}

type dynamicTypeResolver struct {
	*chains.Conn
}

var (
	_ protoregistry.MessageTypeResolver   = dynamicTypeResolver{}
	_ protoregistry.ExtensionTypeResolver = dynamicTypeResolver{}
)

func (d dynamicTypeResolver) FindMessageByName(message protoreflect.FullName) (protoreflect.MessageType, error) {
	desc, err := d.ProtoFiles.FindDescriptorByName(message)
	if err != nil {
		return nil, err
	}

	return dynamicpb.NewMessageType(desc.(protoreflect.MessageDescriptor)), nil
}

func (d dynamicTypeResolver) FindMessageByURL(url string) (protoreflect.MessageType, error) {
	if i := strings.LastIndexByte(url, '/'); i >= 0 {
		url = url[i+len("/"):]
	}

	return d.FindMessageByName(protoreflect.FullName(url))
}

func (d dynamicTypeResolver) FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error) {
	desc, err := d.ProtoFiles.FindDescriptorByName(field)
	if err != nil {
		return nil, err
	}

	return dynamicpb.NewExtensionType(desc.(protoreflect.ExtensionTypeDescriptor)), nil
}

func (d dynamicTypeResolver) FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error) {
	desc, err := d.ProtoFiles.FindDescriptorByName(message)
	if err != nil {
		return nil, err
	}

	messageDesc := desc.(protoreflect.MessageDescriptor)
	exts := messageDesc.Extensions()
	n := exts.Len()
	for i := 0; i < n; i++ {
		ext := exts.Get(i)
		if ext.Number() == field {
			return dynamicpb.NewExtensionType(ext), nil
		}
	}

	return nil, protoregistry.NotFound
}
