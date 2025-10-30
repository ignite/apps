package cmd

import (
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

const (
	flagVersion      = "fee-abstraction-version"
	flagFeeAbsModule = "fee-abstraction"
	flagNoModule     = "no-module"
	flagPath         = "path"

	statusAdding = "Adding Fee Abstraction Module..."

	defaultFeeAbsVersion = "v8.0.2"
	feeAbsModuleName     = "feeabs"

	ScaffoldChainHook  = "scaffold-chain"
	ScaffoldModuleHook = "scaffold-module"
)

// GetCommands returns the list of fee-abstraction app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "fee-abstraction",
			Short: "Integrate the fee abstraction module from osmosis-labs",
			Long:  "Integrate the fee abstraction module from osmosis-labs to make it easy for new chains to accept the currencies of existing chains",
		},
	}
}

// GetHooks returns the list of fee-abstraction app hooks.
func GetHooks() []*plugin.Hook {
	return []*plugin.Hook{
		{
			Name:        ScaffoldChainHook,
			PlaceHookOn: "ignite scaffold chain",
			Flags: []*plugin.Flag{
				{
					Name:         flagFeeAbsModule,
					Usage:        "Create a project that includes the fee abstraction module",
					DefaultValue: "false",
					Type:         plugin.FlagTypeBool,
				},

				{
					Name:         flagVersion,
					Usage:        "fee abstraction semantic version",
					DefaultValue: defaultFeeAbsVersion,
					Type:         plugin.FlagTypeString,
				},
			},
		},
		{
			Name:        ScaffoldModuleHook,
			PlaceHookOn: "ignite scaffold module",
		},
	}
}

func getPath(flags plugin.Flags) string {
	path, _ := flags.GetString(flagPath)
	return path
}

func getVersion(flags plugin.Flags) string {
	version, _ := flags.GetString(flagVersion)
	version = strings.Replace(version, "v", "", 1)
	return version
}

// newChain create new *chain.Chain with home and path flags.
func newChain(chainFolder string, flags plugin.Flags, chainOption ...chain.Option) (*chain.Chain, error) {
	appPath := getPath(flags)
	absPath, err := filepath.Abs(appPath)
	if err != nil {
		return nil, err
	}
	absPath = filepath.Join(absPath, chainFolder)
	return chain.New(absPath, chainOption...)
}
