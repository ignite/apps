package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
	flagPort = "port"
)

type pingPubConfig struct {
	ChainName        string   `json:"chain_name"`
	API              []string `json:"api"`
	RPC              []string `json:"rpc"`
	Coingecko        string   `json:"coingecko"`
	SnapshotProvider string   `json:"snapshot_provider"`
	SdkVersion       string   `json:"sdk_version"`
	CoinType         string   `json:"coin_type"`
	MinTxFee         string   `json:"min_tx_fee"`
	AddrPrefix       string   `json:"addr_prefix"`
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

	// get port
	port, err := flags.GetUint(flagPort)
	if err != nil {
		return errors.Errorf("could not get --%s flag: %s", flagPort, err)
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
		return serve(session, pingPubPath, port)
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
	chainDir := filepath.Join(pingPubPath, "chains", "testnet")
	if err := os.MkdirAll(chainDir, 0755); err != nil {
		return errors.Errorf("failed to create directory %s: %w", chainDir, err)
	}

	// create ping.pub configuration file
	pingCfg := pingPubConfig{
		ChainName:        c.Name(),
		API:              []string{"http://localhost:1317"},
		RPC:              []string{"http://localhost:26657"},
		Coingecko:        "",
		SnapshotProvider: "",
		SdkVersion:       c.Version.String(),
		CoinType:         "118",
		MinTxFee:         "500",
		AddrPrefix:       "cosmos",
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
	_ = session.Printf("🚀 Optionally edit the configuration at %s\n", configFilePath)

	return serve(session, pingPubPath, port)
}

func serve(session *cliui.Session, path string, port uint) error {
	_ = session.Printf("Serving ping.pub explorer at http://localhost:%d\n", port)

	// Validate that path is a directory
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("'%s' is not a valid directory", path)
	}

	// Serve directory content via HTTP
	http.Handle("/", http.FileServer(http.Dir(path)))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
