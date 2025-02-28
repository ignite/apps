package autocli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	addresscodec "cosmossdk.io/core/address"
	"github.com/ignite/apps/connect/internal/autocli/flag"
	"github.com/ignite/apps/connect/internal/flags"
	"github.com/ignite/apps/connect/internal/governance"
	"github.com/ignite/apps/connect/internal/tx"
	"github.com/ignite/apps/connect/internal/util"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// BuildMsgCommand builds the msg commands for all the provided modules.
func (b *Builder) BuildMsgCommand(ctx context.Context, moduleOptions map[string]*autocliv1.ModuleOptions) (*cobra.Command, error) {
	msgCmd := topLevelCmd(ctx, "tx", "Transaction subcommands")

	if err := b.enhanceCommandCommon(msgCmd, msgCmdType, moduleOptions); err != nil {
		return nil, err
	}

	return msgCmd, nil
}

// AddMsgServiceCommands adds a sub-command to the provided command for each
// method in the specified service and returns the command. This can be used in
// order to add auto-generated commands to an existing command.
func (b *Builder) AddMsgServiceCommands(cmd *cobra.Command, cmdDescriptor *autocliv1.ServiceCommandDescriptor) error {
	for cmdName, subCmdDescriptor := range cmdDescriptor.SubCommands {
		subCmd := findSubCommand(cmd, cmdName)
		if subCmd == nil {
			short := subCmdDescriptor.Short
			if short == "" {
				short = fmt.Sprintf("Tx commands for the %s service", subCmdDescriptor.Service)
			}
			subCmd = topLevelCmd(cmd.Context(), cmdName, short)
		}

		// Add recursive sub-commands if there are any. This is used for nested services.
		if err := b.AddMsgServiceCommands(subCmd, subCmdDescriptor); err != nil {
			return err
		}

		if !subCmdDescriptor.EnhanceCustomCommand {
			cmd.AddCommand(subCmd)
		}
	}

	if cmdDescriptor.Service == "" {
		// skip empty command descriptor
		return nil
	}

	descriptor, err := b.FileResolver.FindDescriptorByName(protoreflect.FullName(cmdDescriptor.Service))
	if err != nil {
		return fmt.Errorf("can't find service %s: %w", cmdDescriptor.Service, err)
	}
	service := descriptor.(protoreflect.ServiceDescriptor)
	methods := service.Methods()

	rpcOptMap := map[protoreflect.Name]*autocliv1.RpcCommandOptions{}
	for _, option := range cmdDescriptor.RpcCommandOptions {
		methodName := protoreflect.Name(option.RpcMethod)
		// validate that methods exist
		if m := methods.ByName(methodName); m == nil {
			return fmt.Errorf("rpc method %q not found for service %q", methodName, service.FullName())
		}
		rpcOptMap[methodName] = option

	}

	for i := 0; i < methods.Len(); i++ {
		methodDescriptor := methods.Get(i)
		methodOpts, ok := rpcOptMap[methodDescriptor.Name()]
		if !ok {
			methodOpts = &autocliv1.RpcCommandOptions{}
		}

		if methodOpts.Skip {
			continue
		}

		if !util.IsSupportedVersion(methodDescriptor) {
			continue
		}

		methodCmd, err := b.BuildMsgMethodCommand(methodDescriptor, methodOpts)
		if err != nil {
			return err
		}

		if findSubCommand(cmd, methodCmd.Name()) != nil {
			// do not overwrite existing commands
			// we do not display a warning because you may want to overwrite an autocli command
			continue
		}

		cmd.AddCommand(methodCmd)
	}

	return nil
}

// BuildMsgMethodCommand returns a command that outputs the JSON representation of the message.
func (b *Builder) BuildMsgMethodCommand(descriptor protoreflect.MethodDescriptor, options *autocliv1.RpcCommandOptions) (*cobra.Command, error) {
	execFunc := func(cmd *cobra.Command, input protoreflect.Message) error {
		fd := input.Descriptor().Fields().ByName(protoreflect.Name(flag.GetSignerFieldName(input.Descriptor())))
		addressCodec := b.Builder.AddressCodec

		// handle gov proposals commands
		skipProposal, _ := cmd.Flags().GetBool(flags.FlagNoProposal)
		if isProposalMessage(descriptor.Input()) && !skipProposal {
			return b.handleGovProposal(cmd, input, addressCodec, fd)
		}

		// set signer to signer field if empty
		if addr := input.Get(fd).String(); addr == "" {
			scalarType, ok := flag.GetScalarType(fd)
			if ok {
				// override address codec if validator or consensus address
				switch scalarType {
				case flag.ValidatorAddressStringScalarType:
					addressCodec = b.Builder.ValidatorAddressCodec
				case flag.ConsensusAddressStringScalarType:
					addressCodec = b.Builder.ConsensusAddressCodec
				}
			}

			signer, err := tx.GetFromAddress(cmd)
			if err != nil {
				return fmt.Errorf("failed to get from address: %w", err)
			}

			input.Set(fd, protoreflect.ValueOfString(signer))
		}

		// AutoCLI uses protov2 messages, while the SDK only supports proto v1 messages.
		// Here we use dynamicpb, to create a proto v1 compatible message.
		// The SDK codec will handle protov2 -> protov1 (marshal)
		msg := dynamicpb.NewMessage(input.Descriptor())
		proto.Merge(msg, input.Interface())

		out, err := tx.GenerateAndBroadcastTxCLI(cmd.Context(), nil /* TODO */, msg)
		if err != nil {
			return err
		}

		return b.outOrStdoutFormat(cmd, out)
	}

	cmd, err := b.buildMethodCommandCommon(descriptor, options, execFunc)
	if err != nil {
		return nil, err
	}

	if b.AddTxConnFlags != nil {
		b.AddTxConnFlags(cmd)
	}

	// silence usage only for inner txs & queries commands
	cmd.SilenceUsage = true

	// set gov proposal flags if command is a gov proposal
	if isProposalMessage(descriptor.Input()) {
		governance.AddGovPropFlagsToCmd(cmd)
		cmd.Flags().Bool(flags.FlagNoProposal, false, "Skip gov proposal and submit a normal transaction")
	}

	return cmd, nil
}

// handleGovProposal sets the authority field of the message to the gov module address and creates a gov proposal.
func (b *Builder) handleGovProposal(
	cmd *cobra.Command,
	input protoreflect.Message,
	addressCodec addresscodec.Codec,
	fd protoreflect.FieldDescriptor,
) error {
	govAuthority := authtypes.NewModuleAddress(governance.ModuleName)
	authority, err := addressCodec.BytesToString(govAuthority.Bytes())
	if err != nil {
		return fmt.Errorf("failed to convert gov authority: %w", err)
	}
	input.Set(fd, protoreflect.ValueOfString(authority))

	signerFromFlag, err := tx.GetFromAddress(cmd)
	if err != nil {
		return fmt.Errorf("failed to get from address: %w", err)
	}

	proposal, err := governance.ReadGovPropCmdFlags(signerFromFlag, cmd.Flags())
	if err != nil {
		return err
	}

	// AutoCLI uses protov2 messages, while the SDK only supports proto v1 messages.
	// Here we use dynamicpb, to create a proto v1 compatible message.
	// The SDK codec will handle protov2 -> protov1 (marshal)
	msg := dynamicpb.NewMessage(input.Descriptor())
	proto.Merge(msg, input.Interface())

	if err := governance.SetGovMsgs(proposal, msg); err != nil {
		return fmt.Errorf("failed to set msg in proposal %w", err)
	}

	out, err := tx.GenerateAndBroadcastTxCLI(cmd.Context(), nil /* TODO */, proposal)
	if err != nil {
		return err
	}

	return b.outOrStdoutFormat(cmd, out)
}

// isProposalMessage checks the msg name against well known proposal messages.
// this isn't exhaustive. to have it better we need to add a field in autocli proto
// as it was done in v0.52.
func isProposalMessage(_ protoreflect.MessageDescriptor) bool {
	// msg := []string{
	// 	"cosmos.gov.v1.MsgSubmitProposal",
	// 	"cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
	// 	"cosmos.upgrade.v1beta1.MsgCancelUpgrade",
	// 	"cosmos.distribution.v1beta1.MsgFundCommunityPool",
	// 	"ibc.core.client.v1.MsgIBCSoftwareUpgrade",
	// 	".MsgUpdateParams",
	// }

	// for _, m := range msg {
	// 	if strings.HasSuffix(string(desc.FullName()), m) {
	// 		return true
	// 	}
	// }

	return false
}
