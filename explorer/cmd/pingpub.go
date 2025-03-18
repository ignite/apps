package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/xgit"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

const (
	pingPubGitRepo = "https://github.com/ping-pub/explorer"

	statusCloning     = "Cloning ping.pub explorer..."
	statusConfiguring = "Configuring ping.pub..."

	flagPath = "path"
)

type pingPubConfig struct {
	ChainName  string               `json:"chain_name"`
	API        []pingPubConfigAPI   `json:"api"`
	RPC        []pingPubConfigAPI   `json:"rpc"`
	SdkVersion string               `json:"sdk_version"`
	CoinType   string               `json:"coin_type"`
	MinTxFee   string               `json:"min_tx_fee"`
	Assets     []pingPubConfigAsset `json:"assets"`
	AddrPrefix string               `json:"addr_prefix"`
	ThemeColor string               `json:"theme_color"`
	Logo       string               `json:"logo"`
}

type pingPubConfigAPI struct {
	Address  string `json:"address"`
	Provider string `json:"provider"`
}

type pingPubConfigAsset struct {
	Base        string `json:"base"`
	Symbol      string `json:"symbol"`
	Exponent    string `json:"exponent"`
	CoingeckoID string `json:"coingecko_id"`
	Logo        string `json:"logo"`
}

// ExecutePingPub executes explorer pingpub subcommand.
func ExecutePingPub(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	flags := plugin.Flags(cmd.Flags)

	session := cliui.New(cliui.StartSpinnerWithText(statusCloning))
	defer session.End()

	// get the app path
	appPath, err := flags.GetString(flagPath)
	if err != nil {
		return errors.Errorf("could not get --%s flag: %s", flagPath, err)
	}

	absPath, err := filepath.Abs(appPath)
	if err != nil {
		return err
	}

	// initialize chain object
	c, err := chain.New(absPath, chain.CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// prepare ping.pub directory
	pingPubPath := filepath.Join(absPath, "explorer", "ping-pub")
	if _, err := os.Stat(pingPubPath); err == nil {
		// ping.pub directory already exists, serve it
		return serve(session, pingPubPath)
	}

	// clone ping.pub repository
	if err := xgit.Clone(ctx, pingPubGitRepo, pingPubPath); err != nil {
		return errors.Errorf("failed to clone ping.pub repository: %w", err)
	}

	// remove specified directories and files
	dirsToRemove := []string{
		filepath.Join(pingPubPath, "chains", "mainnet"),
		filepath.Join(pingPubPath, "chains", "testnet"),
	}
	filesToRemove := []string{
		filepath.Join(pingPubPath, "README.md"),
		filepath.Join(pingPubPath, "installation.md"),
	}

	// remove directories
	for _, dir := range dirsToRemove {
		if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
			return errors.Errorf("failed to remove directory %s: %w", dir, err)
		}
	}

	// remove files
	for _, file := range filesToRemove {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			return errors.Errorf("failed to remove file %s: %w", file, err)
		}
	}

	session.StartSpinner(statusConfiguring)

	// create chain directory
	chainDir := filepath.Join(pingPubPath, "chains", "mainnet")
	if err := os.MkdirAll(chainDir, 0755); err != nil {
		return errors.Errorf("failed to create directory %s: %w", chainDir, err)
	}

	chainCfg, err := c.Config()
	if err != nil {
		return errors.Errorf("failed to get chain configuration: %w", err)
	}

	// get validators from config and parse their coins
	// we can assume it holds the base denom
	defaultDenom := "stake"
	if len(chainCfg.Validators) > 0 {
		coin, err := sdk.ParseCoinNormalized(chainCfg.Validators[0].Bonded)
		if err == nil {
			defaultDenom = coin.Denom
		}
	}

	// get bech32 prefix
	bech32Prefix, err := c.Bech32Prefix()
	if err != nil {
		return errors.Errorf("failed to get bech32 prefix: %w", err)
	}

	// get coin type
	coinType, err := c.CoinType()
	if err != nil {
		return errors.Errorf("failed to get coin type: %w", err)
	}

	// create ping.pub configuration file
	pingCfg := pingPubConfig{
		ChainName: c.Name(),
		API: []pingPubConfigAPI{
			{
				Address:  "http://localhost:1317",
				Provider: "localhost",
			},
		},
		RPC: []pingPubConfigAPI{
			{
				Address:  "http://localhost:26657",
				Provider: "localhost",
			},
		},
		Assets: []pingPubConfigAsset{
			{
				Base:     defaultDenom,
				Symbol:   defaultDenom,
				Exponent: "6",
			},
		},
		SdkVersion: c.Version.String(),
		CoinType:   fmt.Sprintf("%d", coinType),
		MinTxFee:   "500",
		AddrPrefix: bech32Prefix,
		ThemeColor: "#467dff",
		Logo:       "/logos/cosmos.svg",
	}
	pingCfgBz, err := json.Marshal(pingCfg)
	if err != nil {
		return errors.Errorf("failed to marshal ping.pub configuration: %w", err)
	}

	configFilePath := filepath.Join(chainDir, fmt.Sprintf("%s.json", c.Name()))
	if err := os.WriteFile(configFilePath, pingCfgBz, 0644); err != nil {
		return errors.Errorf("failed to write ping.pub configuration: %w", err)
	}

	_ = session.Printf("🎉 ping.pub explorer configured successfully at `%s/web/ping-pub`.\n", c.AppPath())
	_ = session.Printf("Optionally edit the configuration at %s\n", configFilePath)

	return serve(session, pingPubPath)
}

func serve(session *cliui.Session, path string) error {
	_ = session.Printf("🚀 Starting ping.pub explorer...\n")

	// check if yarn is installed
	if _, err := exec.LookPath("yarn"); err != nil {
		return errors.New("yarn is not installed. Please install yarn to run the web explorer")
	}

	// run the ping.pub explorer
	cmd := exec.Command("sh", "-c", "yarn --ignore-engines && yarn serve")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return errors.Errorf("failed to start ping.pub explorer: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return errors.Errorf("ping.pub explorer stopped with error: %w", err)
	}

	return nil
}
