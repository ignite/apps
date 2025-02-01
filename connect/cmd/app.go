package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/client/v2/autocli/flag"
	"github.com/cosmos/cosmos-sdk/client"
	sdkflags "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"

	"github.com/ignite/apps/connect/chains"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

func AppHandler(ctx context.Context, cmd *plugin.ExecutedCommand, name string, cfg *chains.ChainConfig, args ...string) (*cobra.Command, error) {
	chainCmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Commands for %s chain", name),
	}

	conn, err := chains.NewConn(name, cfg)
	if err != nil {
		return nil, err
	}

	if err := conn.Load(ctx); err != nil {
		return nil, err
	}

	// add comet commands
	cometCmds := cmtservice.NewCometBFTCommands()
	conn.ModuleOptions[cometCmds.Name()] = cometCmds.AutoCLIOptions()

	appOpts := autocli.AppOptions{
		ModuleOptions: conn.ModuleOptions,
	}

	builder := &autocli.Builder{
		Builder: flag.Builder{
			TypeResolver:          &dynamicTypeResolver{conn},
			FileResolver:          conn.ProtoFiles,
			AddressCodec:          addresscodec.NewBech32Codec(cfg.Bech32Prefix),
			ValidatorAddressCodec: addresscodec.NewBech32Codec(fmt.Sprintf("%svaloper", cfg.Bech32Prefix)),
			ConsensusAddressCodec: addresscodec.NewBech32Codec(fmt.Sprintf("%svalcons", cfg.Bech32Prefix)),
		},
		GetClientConn: func(command *cobra.Command) (grpc.ClientConnInterface, error) {
			return conn.Connect()
		},
		AddQueryConnFlags: func(command *cobra.Command) {
			sdkflags.AddQueryFlagsToCmd(command)
			sdkflags.AddKeyringFlags(command.Flags())
		},
		AddTxConnFlags: sdkflags.AddTxFlagsToCmd,
	}

	// add client context
	clientCtx := client.Context{}
	chainCmd.SetContext(context.WithValue(context.Background(), client.ClientContextKey, &clientCtx))
	if err := appOpts.EnhanceRootCommandWithBuilder(chainCmd, builder); err != nil {
		return nil, err
	}

	if len(args) > 0 {
		chainCmd.SetArgs(args)
	}

	return chainCmd, nil
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
