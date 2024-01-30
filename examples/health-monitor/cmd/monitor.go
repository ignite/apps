package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ignite/cli/ignite/pkg/xurl"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/inancgumus/screen"
	"github.com/spf13/pflag"
)

// ExecuteMonitor executes the monitor subcommand.
func ExecuteMonitor(ctx context.Context, cmd *plugin.ExecutedCommand, chainInfo *plugin.ChainInfo) error {
	flags, err := cmd.NewFlags()
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	jsonFlag, err := getJsonFlag(flags)
	if err != nil {
		return fmt.Errorf("failed to get json flag: %w", err)
	}
	refreshDur, err := getRefreshDurationFlag(flags)
	if err != nil {
		return fmt.Errorf("failed to get refresh-duration flag: %w", err)
	}
	rpcAddress, err := getRpcAddressFlag(flags)
	if err != nil {
		return fmt.Errorf("failed to get rpc-address flag: %w", err)
	}
	if rpcAddress == "" {
		rpcAddress = chainInfo.RpcAddress
	}
	rpcURL, err := xurl.HTTP(rpcAddress)
	if err != nil {
		return fmt.Errorf("invalid rpc address %s: %w", &rpcAddress, err)
	}

	httpClient, err := client.NewClientFromNode(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	ticker := time.NewTicker(refreshDur)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			status, err := httpClient.Status(ctx)
			if err != nil {
				return fmt.Errorf("failed to get status: %w", err)
			}
			if jsonFlag {
				printJson(status)
			} else {
				printUserFriendly(status)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func getJsonFlag(flags *pflag.FlagSet) (bool, error) {
	j, err := flags.GetBool("json")
	if err != nil {
		return false, err
	}
	return j, nil
}

func getRefreshDurationFlag(flags *pflag.FlagSet) (time.Duration, error) {
	r, err := flags.GetString("refresh-duration")
	if err != nil {
		return 0, err
	}
	if r == "" {
		return time.Second * 5, nil
	}
	return time.ParseDuration(r)
}

func getRpcAddressFlag(flags *pflag.FlagSet) (string, error) {
	return flags.GetString("rpc-address")
}

type statusResponse struct {
	Time            time.Time `json:"time"`
	ChainID         string    `json:"chain_id"`
	Version         string    `json:"version"`
	Height          int64     `json:"height"`
	LatestBlockHash string    `json:"latest_block_hash"`
}

func printJson(status *ctypes.ResultStatus) {
	resp := statusResponse{
		Time:            time.Now(),
		ChainID:         status.NodeInfo.Network,
		Version:         status.NodeInfo.Version,
		Height:          status.SyncInfo.LatestBlockHeight,
		LatestBlockHash: status.SyncInfo.LatestBlockHash.String(),
	}
	data, _ := json.Marshal(resp)
	fmt.Println(string(data))
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
