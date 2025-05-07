package network

import (
	"context"
	"strconv"

	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/v28/ignite/pkg/events"
	"github.com/ignite/cli/v28/ignite/pkg/xtime"
	launchtypes "github.com/ignite/network/x/launch/types"
	monitoringctypes "github.com/ignite/network/x/monitoringc/types"
	profiletypes "github.com/ignite/network/x/profile/types"
	projecttypes "github.com/ignite/network/x/project/types"
	rewardtypes "github.com/ignite/network/x/reward/types"
	"github.com/pkg/errors"

	"github.com/ignite/apps/network/network/networktypes"
)

//go:generate mockery --name CosmosClient
type CosmosClient interface {
	Context() client.Context
	Status(ctx context.Context) (*ctypes.ResultStatus, error)
	ConsensusInfo(ctx context.Context, height int64) (cosmosclient.ConsensusInfo, error)
	BroadcastTx(ctx context.Context, account cosmosaccount.Account, msgs ...sdktypes.Msg) (cosmosclient.Response, error)
}

// Network is network builder.
type Network struct {
	node                    Node
	ev                      events.Bus
	cosmos                  CosmosClient
	account                 cosmosaccount.Account
	projectQuery            projecttypes.QueryClient
	launchQuery             launchtypes.QueryClient
	profileQuery            profiletypes.QueryClient
	rewardQuery             rewardtypes.QueryClient
	stakingQuery            stakingtypes.QueryClient
	bankQuery               banktypes.QueryClient
	monitoringConsumerQuery monitoringctypes.QueryClient
	clock                   xtime.Clock
}

//go:generate mockery --name Chain
type Chain interface {
	ID() (string, error)
	ChainID() (string, error)
	Name() string
	SourceURL() string
	SourceHash() string
	GenesisPath() (string, error)
	GentxsPath() (string, error)
	DefaultGentxPath() (string, error)
	AppTOMLPath() (string, error)
	ConfigTOMLPath() (string, error)
	NodeID(ctx context.Context) (string, error)
	CacheBinary(launchID uint64) error
}

type Option func(*Network)

func WithProjectQueryClient(client projecttypes.QueryClient) Option {
	return func(n *Network) {
		n.projectQuery = client
	}
}

func WithProfileQueryClient(client profiletypes.QueryClient) Option {
	return func(n *Network) {
		n.profileQuery = client
	}
}

func WithLaunchQueryClient(client launchtypes.QueryClient) Option {
	return func(n *Network) {
		n.launchQuery = client
	}
}

func WithRewardQueryClient(client rewardtypes.QueryClient) Option {
	return func(n *Network) {
		n.rewardQuery = client
	}
}

func WithStakingQueryClient(client stakingtypes.QueryClient) Option {
	return func(n *Network) {
		n.node.stakingQuery = client
	}
}

func WithMonitoringConsumerQueryClient(client monitoringctypes.QueryClient) Option {
	return func(n *Network) {
		n.monitoringConsumerQuery = client
	}
}

func WithBankQueryClient(client banktypes.QueryClient) Option {
	return func(n *Network) {
		n.bankQuery = client
	}
}

func WithCustomClock(clock xtime.Clock) Option {
	return func(n *Network) {
		n.clock = clock
	}
}

// CollectEvents collects events from the network builder.
func CollectEvents(ev events.Bus) Option {
	return func(n *Network) {
		n.ev = ev
	}
}

// New creates a Builder.
func New(cosmos CosmosClient, account cosmosaccount.Account, options ...Option) Network {
	n := Network{
		cosmos:                  cosmos,
		account:                 account,
		node:                    NewNode(cosmos),
		projectQuery:            projecttypes.NewQueryClient(cosmos.Context()),
		launchQuery:             launchtypes.NewQueryClient(cosmos.Context()),
		profileQuery:            profiletypes.NewQueryClient(cosmos.Context()),
		rewardQuery:             rewardtypes.NewQueryClient(cosmos.Context()),
		stakingQuery:            stakingtypes.NewQueryClient(cosmos.Context()),
		bankQuery:               banktypes.NewQueryClient(cosmos.Context()),
		monitoringConsumerQuery: monitoringctypes.NewQueryClient(cosmos.Context()),
		clock:                   xtime.NewClockSystem(),
	}
	for _, opt := range options {
		opt(&n)
	}
	return n
}

func ParseID(id string) (uint64, error) {
	objID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "error parsing ID")
	}
	return objID, nil
}

// AccountAddress returns the address of the account used by the network builder.
func (n Network) AccountAddress() (string, error) {
	return n.account.Address(networktypes.SPN)
}
