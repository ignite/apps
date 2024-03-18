package cmd

import (
	"context"
	"fmt"

	cowsay "github.com/Code-Hex/Neo-cowsay/v2"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

// ExecuteCowsay executes the cowsay subcommand.
func ExecuteCowsay(_ context.Context, cmd *plugin.ExecutedCommand) error {
	flags, err := cmd.NewFlags()
	if err != nil {
		return err
	}

	var (
		name, _ = flags.GetString("name")
		typ, _  = flags.GetString("type")
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
