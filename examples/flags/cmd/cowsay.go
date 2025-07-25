package cmd

import (
	"context"
	"fmt"

	cowsay "github.com/Code-Hex/Neo-cowsay/v2"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

// ExecuteCowsay executes the cowsay subcommand.
func ExecuteCowsay(_ context.Context, cmd *plugin.ExecutedCommand) error {
	var (
		flags   = plugin.Flags(cmd.Flags)
		name, _ = flags.GetString(flagName)
		typ, _  = flags.GetString(flagType)
	)
	say, err := cowsay.Say(
		fmt.Sprintf("Hello, %s!", name),
		cowsay.Type(typ),
		cowsay.BallonWidth(40),
	)
	if err != nil {
		return errors.Errorf("internal error with cowsay: %s", err)
	}
	fmt.Println(say)
	return nil
}
