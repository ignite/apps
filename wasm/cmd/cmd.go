package cmd

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/wasm/services/scaffolder"
)

const (
	flagVersion = "version"

	statusScaffolding  = "Scaffolding..."
	statusAddingConfig = "Adding config..."
)

var cfgFlags = []*plugin.Flag{
	{
		Name:         flagSimulationGasLimit,
		Usage:        "the max gas to be used in a tx simulation call. When not set the consensus max block gas is used instead",
		DefaultValue: "0",
		Type:         plugin.FlagTypeUint64,
	},
	{
		Name:         flagSmartQueryGasLimit,
		Usage:        "the max gas to be used in a smart query contract call",
		DefaultValue: "3000000",
		Type:         plugin.FlagTypeUint64,
	},
	{
		Name:         flagMemoryCacheSize,
		Usage:        "memory cache size in MiB not bytes",
		DefaultValue: "100",
		Type:         plugin.FlagTypeUint64,
	},
}

// GetCommands returns the list of extension commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:     "wasm [command]",
			Aliases: []string{"w"},
			Short:   "Ignite wasm integration",
			Commands: []*plugin.Command{
				{
					Use:   "add",
					Short: "Add wasm support",
					Flags: append(cfgFlags,
						&plugin.Flag{
							Name:         flagVersion,
							Usage:        "wasmd semantic version",
							Shorthand:    "v",
							DefaultValue: scaffolder.DefaultWasmVersion.String(),
							Type:         plugin.FlagTypeString,
						},
					),
				},
				{
					Use:   "config",
					Short: "Add wasm config support",
					Flags: cfgFlags,
				},
			},
		},
	}
}

var (
	modifyPrefix = colors.Modified("modify ")
	createPrefix = colors.Success("create ")
	removePrefix = func(s string) string {
		return strings.TrimPrefix(strings.TrimPrefix(s, modifyPrefix), createPrefix)
	}
)

func getVersion(flags plugin.Flags) string {
	version, _ := flags.GetString(flagVersion)
	version = strings.Replace(version, "v", "", 1)
	return version
}

func getSimulationGasLimit(flags plugin.Flags) uint64 {
	simulationGasLimit, _ := flags.GetUint64(flagSimulationGasLimit)
	return simulationGasLimit
}

func getSmartQueryGasLimit(flags plugin.Flags) uint64 {
	smartQueryGasLimit, _ := flags.GetUint64(flagSmartQueryGasLimit)
	return smartQueryGasLimit
}

func getMemoryCacheSize(flags plugin.Flags) uint64 {
	memoryCacheSize, _ := flags.GetUint64(flagMemoryCacheSize)
	return memoryCacheSize
}

// newChain create new *chain.Chain with home and path.
func newChain(ctx context.Context, api plugin.ClientAPI, chainOption ...chain.Option) (*chain.Chain, error) {
	info, err := api.GetChainInfo(ctx)
	if err != nil {
		return nil, err
	}

	// Check if a custom home is provided
	if info.Home != "" {
		chainOption = append(chainOption, chain.HomePath(info.Home))
	}

	absPath, err := filepath.Abs(info.AppPath)
	if err != nil {
		return nil, err
	}

	return chain.New(absPath, chainOption...)
}
