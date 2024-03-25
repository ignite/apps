package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/xurl"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/inancgumus/screen"
)

// ExecuteMonitor executes the monitor subcommand.
func ExecuteMonitor(ctx context.Context, cmd *plugin.ExecutedCommand, chainInfo *plugin.ChainInfo) error {
	flags, err := cmd.NewFlags()
	if err != nil {
		return errors.Errorf("failed to parse flags: %s", err)
	}

	var (
		isJSON, _          = flags.GetBool(flagJSON)
		refreshDuration, _ = flags.GetString(flagRefreshDuration)
		rpcAddress, _      = flags.GetString(flagRPCAddress)
	)

	if rpcAddress == "" {
		rpcAddress = chainInfo.RpcAddress
	}
	rpcURL, err := xurl.HTTP(rpcAddress)
	if err != nil {
		return errors.Errorf("invalid rpc address %s: %s", rpcAddress, err)
	}

	// Create a Cosmos client instance
	client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(rpcURL))
	if err != nil {
		return errors.Errorf("failed to create client: %s", err)
	}

	refresh, err := time.ParseDuration(refreshDuration)
	if err != nil {
		return errors.Errorf("failed to parse %s flag: %s", flagRefreshDuration, err)
	}
	ticker := time.NewTicker(refresh)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			status, err := client.Status(ctx)
			if err != nil {
				return errors.Errorf("failed to get status: %s", err)
			}
			if err := printStatus(isJSON, status); err != nil {
				return errors.Errorf("failed to print status: %s", err)
			}
		}
	}
}

func printStatus(isJSON bool, status *ctypes.ResultStatus) error {
	if isJSON {
		return printJSON(status)
	}

	screen.Clear()
	screen.MoveTopLeft()
	fmt.Printf("Time: %s\n", time.Now().Format(time.DateTime))
	fmt.Printf("Chain ID: %s\n", status.NodeInfo.Network)
	fmt.Printf("Version: %s\n", status.NodeInfo.Version)
	fmt.Printf("Height: %d\n", status.SyncInfo.LatestBlockHeight)
	fmt.Printf("Latest Block Hash: %s\n", status.SyncInfo.LatestBlockHash.String())

	return nil
}

type statusResponse struct {
	Time            time.Time `json:"time"`
	ChainID         string    `json:"chain_id"`
	Version         string    `json:"version"`
	Height          int64     `json:"height"`
	LatestBlockHash string    `json:"latest_block_hash"`
}

func printJSON(status *ctypes.ResultStatus) error {
	resp := statusResponse{
		Time:            time.Now(),
		ChainID:         status.NodeInfo.Network,
		Version:         status.NodeInfo.Version,
		Height:          status.SyncInfo.LatestBlockHeight,
		LatestBlockHash: status.SyncInfo.LatestBlockHash.String(),
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
