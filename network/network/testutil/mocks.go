package testutil

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	launchtypes "github.com/ignite/network/x/launch/types"
	monitoringctypes "github.com/ignite/network/x/monitoringc/types"
	monitoringptypes "github.com/ignite/network/x/monitoringp/types"
	profiletypes "github.com/ignite/network/x/profile/types"
	projecttypes "github.com/ignite/network/x/project/types"
	rewardtypes "github.com/ignite/network/x/reward/types"
)

type ProjectClient interface {
	projecttypes.QueryClient
}

type ProfileClient interface {
	profiletypes.QueryClient
}

type LaunchClient interface {
	launchtypes.QueryClient
}

type RewardClient interface {
	rewardtypes.QueryClient
}

type BankClient interface {
	banktypes.QueryClient
}

type StakingClient interface {
	stakingtypes.QueryClient
}

type MonitoringcClient interface {
	monitoringctypes.QueryClient
}

type MonitoringpClient interface {
	monitoringptypes.QueryClient
}

type AccountInfo interface {
	keyring.Record
}
