package cmd

import (
	"context"
	"fmt"

	cowsay "github.com/Code-Hex/Neo-cowsay/v2"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/spf13/pflag"
)

// ExecuteCowsay executes the cowsay subcommand.
func ExecuteCowsay(_ context.Context, cmd *plugin.ExecutedCommand) error {
	flags, err := cmd.NewFlags()
	if err != nil {
		return err
	}

	name, err := getNameFlag(flags)
	if err != nil {
		return err
	}
	typ, err := getTypeflag(flags)
	if err != nil {
		return err
	}
	say, err := cowsay.Say(
		fmt.Sprintf("Hello, %s!", name),
		cowsay.Type(typ),
		cowsay.BallonWidth(40),
	)
	if err != nil {
		return fmt.Errorf("internal error with cowsay: %w", err)
	}
	fmt.Println(say)
	return nil
}

func getTypeflag(flags *pflag.FlagSet) (string, error) {
	typ, err := flags.GetString("type")
	if err != nil {
		return "", fmt.Errorf("could not get --type flag: %w", err)
	}

	if typ == "" {
		return "default", nil
	}
	return typ, nil
}
