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
	"github.com/spf13/pflag"
)

const (
	flagPath    = "path"
	flagHome    = "home"
	flagVersion = "version"

	statusScaffolding  = "Scaffolding..."
	statusAddingConfig = "Adding config..."

	defaultSimulationGasLimit = 0
	defaultSmartQueryGasLimit = 3_000_000
	defaultMemoryCacheSize    = 100
	defaultWasmVersion        = "v0.50.0"
)

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
					Flags: []*plugin.Flag{
						{
							Name:      flagPath,
							Usage:     "path of the app",
							Shorthand: "p",
							Type:      plugin.FlagTypeString,
						},
						{
							Name:  flagHome,
							Usage: "directory where the blockchain node is initialized",
							Type:  plugin.FlagTypeString,
						},
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
						{
							Name:         flagVersion,
							Usage:        "wasmd semantic version",
							Shorthand:    "v",
							DefaultValue: defaultWasmVersion,
							Type:         plugin.FlagTypeString,
						},
					},
				},
				{
					Use:   "config",
					Short: "Add wasm config support",
					Flags: []*plugin.Flag{
						{
							Name:      flagPath,
							Usage:     "path of the app",
							Shorthand: "p",
							Type:      plugin.FlagTypeString,
						},
						{
							Name:  flagHome,
							Usage: "directory where the blockchain node is initialized",
							Type:  plugin.FlagTypeString,
						},
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
					},
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

func getPath(flags *pflag.FlagSet) string {
	path, _ := flags.GetString(flagPath)
	return path
}

func getHome(flags *pflag.FlagSet) string {
	home, _ := flags.GetString(flagHome)
	return home
}

func getWasmVersion(flags *pflag.FlagSet) string {
	version, err := flags.GetString(flagVersion)
	if err != nil || version == "" {
		version = defaultWasmVersion
	}
	version = strings.Replace(version, "v", "", 1)
	return version
}

func getSimulationGasLimit(flags *pflag.FlagSet) uint64 {
	simulationGasLimit, err := flags.GetUint64(flagSimulationGasLimit)
	if err != nil || simulationGasLimit == 0 {
		simulationGasLimit = defaultSimulationGasLimit
	}
	return simulationGasLimit
}

func getSmartQueryGasLimit(flags *pflag.FlagSet) uint64 {
	smartQueryGasLimit, err := flags.GetUint64(flagSmartQueryGasLimit)
	if err != nil || smartQueryGasLimit == 0 {
		smartQueryGasLimit = defaultSmartQueryGasLimit
	}
	return smartQueryGasLimit
}

func getMemoryCacheSize(flags *pflag.FlagSet) uint64 {
	memoryCacheSize, err := flags.GetUint64(flagMemoryCacheSize)
	if err != nil || memoryCacheSize == 0 {
		memoryCacheSize = defaultMemoryCacheSize
	}
	return memoryCacheSize
}

// newChainWithHomeFlags create new *chain.Chain with home and path flags.
func newChainWithHomeFlags(flags *pflag.FlagSet, chainOption ...chain.Option) (*chain.Chain, error) {
	// Check if custom home is provided
	if home := getHome(flags); home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}

	appPath := getPath(flags)
	absPath, err := filepath.Abs(appPath)
	if err != nil {
		return nil, err
	}

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
