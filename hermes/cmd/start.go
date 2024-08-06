package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/hermes/pkg/hermes"
)

func StartHandler(ctx context.Context, cmd *plugin.ExecutedCommand) (err error) {
	var (
		flags     = cmd.Flags
		args      = cmd.Args
		customCfg = getConfig(flags)
		cfgName   = strings.Join(args, hermes.ConfigNameSeparator)
	)

	cfgPath := customCfg
	if cfgPath == "" {
		cfgPath, err = hermes.ConfigPath(cfgName)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return errors.Errorf("config file (%s) not exist, try to configure you relayer first", cfgPath)
	}

	h, err := hermes.New()
	if err != nil {
		return err
	}
	defer h.Cleanup()

	return h.Start(
		ctx,
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdIn(os.Stdin),
		hermes.WithStdOut(os.Stdout),
		hermes.WithStdErr(os.Stderr),
	)
}
