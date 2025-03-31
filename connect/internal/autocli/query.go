package autocli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"cosmossdk.io/math"
	"cosmossdk.io/x/tx/signing/aminojson"
	"github.com/ignite/apps/connect/internal/flags"
	"github.com/ignite/apps/connect/internal/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BuildQueryCommand builds the query commands for all the provided modules.
func (b *Builder) BuildQueryCommand(ctx context.Context, moduleOptions map[string]*autocliv1.ModuleOptions) (*cobra.Command, error) {
	queryCmd := topLevelCmd(ctx, "query", "Querying subcommands")
	queryCmd.Aliases = []string{"q"}

	if err := b.enhanceCommandCommon(queryCmd, queryCmdType, moduleOptions); err != nil {
		return nil, err
	}

	return queryCmd, nil
}

// AddQueryServiceCommands adds a sub-command to the provided command for each
// method in the specified service and returns the command. This can be used in
// order to add auto-generated commands to an existing command.
func (b *Builder) AddQueryServiceCommands(cmd *cobra.Command, cmdDescriptor *autocliv1.ServiceCommandDescriptor) error {
	for cmdName, subCmdDesc := range cmdDescriptor.SubCommands {
		subCmd := findSubCommand(cmd, cmdName)
		if subCmd == nil {
			short := subCmdDesc.Short
			if short == "" {
				short = fmt.Sprintf("Querying commands for the %s service", subCmdDesc.Service)
			}
			subCmd = topLevelCmd(cmd.Context(), cmdName, short)
		}

		if err := b.AddQueryServiceCommands(subCmd, subCmdDesc); err != nil {
			return err
		}

		if !subCmdDesc.EnhanceCustomCommand {
			cmd.AddCommand(subCmd)
		}
	}

	// skip empty command descriptors
	if cmdDescriptor.Service == "" {
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
		name := protoreflect.Name(option.RpcMethod)
		rpcOptMap[name] = option
		// make sure method exists
		if m := methods.ByName(name); m == nil {
			return fmt.Errorf("rpc method %q not found for service %q", name, service.FullName())
		}
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

		methodCmd, err := b.BuildQueryMethodCommand(cmd.Context(), methodDescriptor, methodOpts)
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

// BuildQueryMethodCommand creates a gRPC query command for the given service method. This can be used to auto-generate
// just a single command for a single service rpc method.
func (b *Builder) BuildQueryMethodCommand(ctx context.Context, descriptor protoreflect.MethodDescriptor, options *autocliv1.RpcCommandOptions) (*cobra.Command, error) {
	serviceDescriptor := descriptor.Parent().(protoreflect.ServiceDescriptor)
	methodName := fmt.Sprintf("/%s/%s", serviceDescriptor.FullName(), descriptor.Name())
	outputType := util.ResolveMessageType(b.TypeResolver, descriptor.Output())
	encoderOptions := aminojson.EncoderOptions{
		Indent:             "  ",
		EnumAsString:       true,
		DoNotSortFields:    true,
		AminoNameAsTypeURL: true,
		MarshalMappings:    true,
		TypeResolver:       b.TypeResolver,
		FileResolver:       b.FileResolver,
	}

	cmd, err := b.buildMethodCommandCommon(descriptor, options, func(cmd *cobra.Command, input protoreflect.Message) error {
		clientConn, err := b.GetClientConn(cmd)
		if err != nil {
			return err
		}

		output := outputType.New()
		if err := clientConn.Invoke(b.queryContext(cmd.Context(), cmd), methodName, input.Interface(), output.Interface()); err != nil {
			return err
		}

		if noIndent, _ := cmd.Flags().GetBool(flags.FlagNoIndent); noIndent {
			encoderOptions.Indent = ""
		}

		enc := encoder(aminojson.NewEncoder(encoderOptions))
		bz, err := enc.Marshal(output.Interface())
		if err != nil {
			return fmt.Errorf("cannot marshal response %v: %w", output.Interface(), err)
		}

		return b.outOrStdoutFormat(cmd, bz)
	})
	if err != nil {
		return nil, err
	}

	if b.AddQueryConnFlags != nil {
		b.AddQueryConnFlags(cmd)

		cmd.Flags().BoolP(flags.FlagNoIndent, "", false, "Do not indent JSON output")
	}

	// silence usage only for inner txs & queries commands
	if cmd != nil {
		cmd.SilenceUsage = true
	}

	return cmd, nil
}

// queryContext returns a new context with metadata for block height if specified.
// If the context already has metadata, it is returned as-is. Otherwise, if a height
// flag is present on the command, it adds an x-cosmos-block-height metadata value
// with the specified height.
func (b *Builder) queryContext(ctx context.Context, cmd *cobra.Command) context.Context {
	md, _ := metadata.FromOutgoingContext(ctx)
	if md != nil {
		return ctx
	}

	md = map[string][]string{}
	if cmd.Flags().Lookup(flags.FlagHeight) != nil {
		h, _ := cmd.Flags().GetInt64(flags.FlagHeight)
		md["x-cosmos-block-height"] = []string{fmt.Sprintf("%d", h)}
	}

	return metadata.NewOutgoingContext(ctx, md)
}

func encoder(encoder aminojson.Encoder) aminojson.Encoder {
	return encoder.DefineTypeEncoding("google.protobuf.Duration", func(_ *aminojson.Encoder, msg protoreflect.Message, w io.Writer) error {
		var (
			secondsName protoreflect.Name = "seconds"
			nanosName   protoreflect.Name = "nanos"
		)

		fields := msg.Descriptor().Fields()
		secondsField := fields.ByName(secondsName)
		if secondsField == nil {
			return errors.New("expected seconds field")
		}

		seconds := msg.Get(secondsField).Int()

		nanosField := fields.ByName(nanosName)
		if nanosField == nil {
			return errors.New("expected nanos field")
		}

		nanos := msg.Get(nanosField).Int()

		_, err := fmt.Fprintf(w, `"%s"`, (time.Duration(seconds)*time.Second + (time.Duration(nanos) * time.Nanosecond)).String())
		return err
	}).DefineTypeEncoding("cosmos.base.v1beta1.DecCoin", func(_ *aminojson.Encoder, msg protoreflect.Message, w io.Writer) error {
		var (
			denomName  protoreflect.Name = "denom"
			amountName protoreflect.Name = "amount"
		)

		fields := msg.Descriptor().Fields()
		denomField := fields.ByName(denomName)
		if denomField == nil {
			return errors.New("expected denom field")
		}

		denom := msg.Get(denomField).String()

		amountField := fields.ByName(amountName)
		if amountField == nil {
			return errors.New("expected amount field")
		}

		amount := msg.Get(amountField).String()
		decimalPlace := len(amount) - math.LegacyPrecision
		if decimalPlace > 0 {
			amount = amount[:decimalPlace] + "." + amount[decimalPlace:]
		} else if decimalPlace == 0 {
			amount = "0." + amount
		} else {
			amount = "0." + strings.Repeat("0", -decimalPlace) + amount
		}

		amountDec, err := math.LegacyNewDecFromStr(amount)
		if err != nil {
			return fmt.Errorf("invalid amount: %s: %w", amount, err)
		}

		_, err = fmt.Fprintf(w, `"%s"`, sdk.NewDecCoinFromDec(denom, amountDec))
		return err
	})
}
