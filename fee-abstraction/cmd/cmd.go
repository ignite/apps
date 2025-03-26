package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

const (
	flagVersion      = "fee-abstraction-version"
	flagFeeAbsModule = "fee-abstraction"
	flagNoModule     = "no-module"
	flagPath         = "path"

	statusScaffolding = "Scaffolding..."

	defaultFeeAbsVersion = "v8.0.2"
	feeAbsModuleName     = "feeabs"

	ScaffoldChainHook  = "scaffold-chain"
	ScaffoldModuleHook = "scaffold-module"
)

var (
	modifyPrefix = colors.Modified("modify ")
	createPrefix = colors.Success("create ")
	removePrefix = func(s string) string {
		return strings.TrimPrefix(strings.TrimPrefix(s, modifyPrefix), createPrefix)
	}
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
					Shorthand:    "v",
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

// sourceModificationToString output the modifications into a readable text.
func sourceModificationToString(sm xgenny.SourceModification) (string, error) {
	// get file names and add prefix
	var files []string
	for _, modified := range sm.ModifiedFiles() {
		// get the relative app path from the current directory
		relativePath, err := relativePath(modified)
		if err != nil {
			return "", err
		}
		files = append(files, modifyPrefix+relativePath)
	}
	for _, created := range sm.CreatedFiles() {
		// get the relative app path from the current directory
		relativePath, err := relativePath(created)
		if err != nil {
			return "", err
		}
		files = append(files, createPrefix+relativePath)
	}

	// sort filenames without prefix
	sort.Slice(files, func(i, j int) bool {
		s1 := removePrefix(files[i])
		s2 := removePrefix(files[j])

		return strings.Compare(s1, s2) == -1
	})

	return "\n" + strings.Join(files, "\n"), nil
}

// relativePath return the relative app path from the current directory.
func relativePath(appPath string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path, err := filepath.Rel(pwd, appPath)
	if err != nil {
		return "", err
	}
	return path, nil
}
