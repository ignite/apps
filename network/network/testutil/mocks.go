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

//go:generate mockery --name ProjectClient
type ProjectClient interface {
	projecttypes.QueryClient
}

//go:generate mockery --name ProfileClient
type ProfileClient interface {
	profiletypes.QueryClient
}

//go:generate mockery --name LaunchClient
type LaunchClient interface {
	launchtypes.QueryClient
}

//go:generate mockery --name RewardClient
type RewardClient interface {
	rewardtypes.QueryClient
}

//go:generate mockery --name BankClient
type BankClient interface {
	banktypes.QueryClient
}

//go:generate mockery --name StakingClient
type StakingClient interface {
	stakingtypes.QueryClient
}

//go:generate mockery --name MonitoringcClient
type MonitoringcClient interface {
	monitoringctypes.QueryClient
}

//go:generate mockery --name MonitoringpClient
type MonitoringpClient interface {
	monitoringptypes.QueryClient
}

//go:generate mockery --name AccountInfo
type AccountInfo interface {
	keyring.Record
}
