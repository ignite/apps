package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/spf13/cobra"
)

// NewWasm creates a new wasm command that holds
// some other sub commands related to CosmWasm.
func NewWasm() *cobra.Command {
	c := &cobra.Command{
		Use:           "wasm [command]",
		Aliases:       []string{"w"},
		Short:         "Ignite wasm integration",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// add sub commands.
	c.AddCommand(
		NewWasmAdd(),
		NewWasmConfig(),
	)
	return c
}

const (
	flagPath    = "path"
	flagHome    = "home"
	flagVersion = "version"

	statusScaffolding  = "Scaffolding..."
	statusAddingConfig = "Adding config..."

	defaultWasmVersion = "v0.50.0"
)

var (
	modifyPrefix = colors.Modified("modify ")
	createPrefix = colors.Success("create ")
	removePrefix = func(s string) string {
		return strings.TrimPrefix(strings.TrimPrefix(s, modifyPrefix), createPrefix)
	}
)

func flagSetPath(cmd *cobra.Command) {
	cmd.Flags().StringP(flagPath, "p", ".", "path of the app")
}

func flagSetHome(cmd *cobra.Command) {
	cmd.Flags().String(flagHome, "", "directory where the blockchain node is initialized")
}

func flagSetWasmVersion(cmd *cobra.Command) {
	cmd.Flags().String(flagVersion, defaultWasmVersion, "wasmd semantic version")
}

func flagSetWasmConfigs(cmd *cobra.Command) {
	cmd.Flags().Uint64(flagSimulationGasLimit, 0, "the max gas to be used in a tx simulation call. When not set the consensus max block gas is used instead")
	cmd.Flags().Uint64(flagSmartQueryGasLimit, 3_000_000, "the max gas to be used in a smart query contract call")
	cmd.Flags().Uint64(flagMemoryCacheSize, 100, "memory cache size in MiB not bytes")
}

func getPath(cmd *cobra.Command) string {
	path, _ := cmd.Flags().GetString(flagPath)
	return path
}

func getHome(cmd *cobra.Command) string {
	home, _ := cmd.Flags().GetString(flagHome)
	return home
}

func getWasmVersion(cmd *cobra.Command) string {
	version, _ := cmd.Flags().GetString(flagVersion)
	version = strings.Replace(version, "v", "", 1)
	return version
}

func getSimulationGasLimit(cmd *cobra.Command) uint64 {
	simulationGasLimit, _ := cmd.Flags().GetUint64(flagSimulationGasLimit)
	return simulationGasLimit
}

func getSmartQueryGasLimit(cmd *cobra.Command) uint64 {
	smartQueryGasLimit, _ := cmd.Flags().GetUint64(flagSmartQueryGasLimit)
	return smartQueryGasLimit
}

func getMemoryCacheSize(cmd *cobra.Command) uint64 {
	memoryCacheSize, _ := cmd.Flags().GetUint64(flagMemoryCacheSize)
	return memoryCacheSize
}

// newChainWithHomeFlags create new *chain.Chain with home and path flags.
func newChainWithHomeFlags(cmd *cobra.Command, chainOption ...chain.Option) (*chain.Chain, error) {
	// Check if custom home is provided
	if home := getHome(cmd); home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}

	appPath := getPath(cmd)
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
