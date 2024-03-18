package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
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
		jsonFlag, _   = flags.GetBool("json")
		refreshDur, _ = flags.GetDuration("refresh-duration")
		rpcAddress, _ = flags.GetString("rpc-address")
	)

	if rpcAddress == "" {
		rpcAddress = chainInfo.RpcAddress
	}
	rpcURL, err := xurl.TCP(rpcAddress)
	if err != nil {
		return errors.Errorf("invalid rpc address %s: %s", rpcAddress, err)
	}

	httpClient, err := client.NewClientFromNode(rpcURL)
	if err != nil {
		return errors.Errorf("failed to create client: %s", err)
	}

	ticker := time.NewTicker(refreshDur)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			status, err := httpClient.Status(ctx)
			if err != nil {
				return errors.Errorf("failed to get status: %s", err)
			}
			if jsonFlag {
				if err := printJSON(status); err != nil {
					return err
				}
			} else {
				printUserFriendly(status)
			}
		}
	}
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

func printUserFriendly(status *ctypes.ResultStatus) {
	screen.Clear()
	screen.MoveTopLeft()
	fmt.Printf("Time: %s\n", time.Now().Format(time.DateTime))
	fmt.Printf("Chain ID: %s\n", status.NodeInfo.Network)
	fmt.Printf("Version: %s\n", status.NodeInfo.Version)
	fmt.Printf("Height: %d\n", status.SyncInfo.LatestBlockHeight)
	fmt.Printf("Latest Block Hash: %s\n", status.SyncInfo.LatestBlockHash.String())
}
