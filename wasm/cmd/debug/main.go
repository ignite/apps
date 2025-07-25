package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/wasm/cmd"
	"github.com/ignite/apps/wasm/services/scaffolder"
)

const (
	flagSimulationGasLimit = "simulation-gas-limit"
	flagSmartQueryGasLimit = "query-gas-limit"
	flagMemoryCacheSize    = "memory-cache-size"
	flagVersion            = "version"
)

type apiMock struct{}

func (m apiMock) GetChainInfo(context.Context) (*plugin.ChainInfo, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}
	return &plugin.ChainInfo{
		ChainId:    "venus",
		AppPath:    wd,
		ConfigPath: wd + "/config.yaml",
		RpcAddress: "http://localhost:26657",
		Home:       wd + "/home",
	}, nil
}

func (m apiMock) GetIgniteInfo(context.Context) (*plugin.IgniteInfo, error) {
	return &plugin.IgniteInfo{
		CliVersion: "v29.0.0",
		GoVersion: func() string {
			if info, ok := os.LookupEnv("GOVERSION"); ok {
				return info
			}
			return "go1.20"
		}(),
		SdkVersion:      "v0.53.2",
		BufVersion:      "",
		BuildDate:       "",
		SourceHash:      "",
		ConfigVersion:   "",
		Os:              "",
		Arch:            "",
		BuildFromSource: false,
	}, nil
}

var cfgFlags = []*plugin.Flag{
	{
		Name:         flagSimulationGasLimit,
		Usage:        "the max gas to be used in a tx simulation call. When not set the consensus max block gas is used instead",
		DefaultValue: "0",
		Value:        "0",
		Type:         plugin.FlagTypeUint64,
	},
	{
		Name:         flagSmartQueryGasLimit,
		Usage:        "the max gas to be used in a smart query contract call",
		DefaultValue: "3000000",
		Value:        "3000000",
		Type:         plugin.FlagTypeUint64,
	},
	{
		Name:         flagMemoryCacheSize,
		Usage:        "memory cache size in MiB not bytes",
		DefaultValue: "100",
		Value:        "100",
		Type:         plugin.FlagTypeUint64,
	},
}

func main() {
	var (
		args    = os.Args[1:]
		ctx     = context.Background()
		api     = apiMock{}
		cmdName = args[0]
		c       = &plugin.ExecutedCommand{
			Use:    cmdName,
			Path:   "ignite wasm " + cmdName,
			Args:   args[1:],
			OsArgs: args,
		}
	)
	switch cmdName {
	case "add":
		c.Flags = append(c.Flags, append(cfgFlags,
			&plugin.Flag{
				Name:         flagVersion,
				Usage:        "wasmd semantic version",
				Shorthand:    "v",
				DefaultValue: scaffolder.DefaultWasmVersion.String(),
				Value:        scaffolder.DefaultWasmVersion.String(),
				Type:         plugin.FlagTypeString,
			})...)
		if err := cmd.AddHandler(ctx, c, api); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	case "config":
		if err := cmd.ConfigHandler(ctx, c, api); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown hermes command: %s", cmdName)
		return
	}
}
