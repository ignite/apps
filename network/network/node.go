package network

import (
	"context"
	"encoding/base64"
	"errors"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcconntypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	spntypes "github.com/ignite/network/pkg/types"
	monitoringptypes "github.com/ignite/network/x/monitoringp/types"

	"github.com/ignite/apps/network/network/networktypes"
)

// Node is node builder.
type Node struct {
	cosmos                  CosmosClient
	stakingQuery            stakingtypes.QueryClient
	ibcClientQuery          ibcclienttypes.QueryClient
	ibcConnQuery            ibcconntypes.QueryClient
	ibcChannelQuery         ibcchanneltypes.QueryClient
	monitoringProviderQuery monitoringptypes.QueryClient
}

// NewNode creates a new client for node API.
func NewNode(cosmos CosmosClient) Node {
	return Node{
		cosmos:                  cosmos,
		stakingQuery:            stakingtypes.NewQueryClient(cosmos.Context()),
		ibcClientQuery:          ibcclienttypes.NewQueryClient(cosmos.Context()),
		ibcConnQuery:            ibcconntypes.NewQueryClient(cosmos.Context()),
		ibcChannelQuery:         ibcchanneltypes.NewQueryClient(cosmos.Context()),
		monitoringProviderQuery: monitoringptypes.NewQueryClient(cosmos.Context()),
	}
}

// RewardsInfo Fetches the consensus state with the validator set and the unbounding time.
func (n Node) RewardsInfo(ctx context.Context) (
	reward networktypes.Reward,
	chainID string,
	unbondingTimeSeconds int64,
	err error,
) {
	status, err := n.cosmos.Status(ctx)
	if err != nil {
		return networktypes.Reward{}, "", 0, err
	}
	lastBlockHeight := status.SyncInfo.LatestBlockHeight

	info, err := n.consensus(ctx, n.cosmos, lastBlockHeight)
	if err != nil {
		return networktypes.Reward{}, "", 0, err
	}

	stakingParams, err := n.stakingParams(ctx)
	if err != nil {
		return networktypes.Reward{}, "", 0, err
	}

	return info,
		status.NodeInfo.Network,
		int64(stakingParams.UnbondingTime.Seconds()),
		nil
}

// RewardIBCInfo returns IBC info to relay packets to claim rewards for the chain.
func (n Node) RewardIBCInfo(ctx context.Context) (networktypes.RewardIBCInfo, error) {
	clientID, err := n.consumerClientID(ctx)
	if err != nil && !errors.Is(err, ErrObjectNotFound) {
		return networktypes.RewardIBCInfo{}, err
	}

	channelID, err := n.connectionChannelID(ctx)
	if err != nil && !errors.Is(err, ErrObjectNotFound) {
		return networktypes.RewardIBCInfo{}, err
	}

	connections, err := n.clientConnections(ctx, clientID)
	if err != nil && !errors.Is(err, ErrObjectNotFound) {
		return networktypes.RewardIBCInfo{}, err
	}

	info := networktypes.RewardIBCInfo{
		ClientID:  clientID,
		ChannelID: channelID,
	}

	if len(connections) > 0 {
		info.ConnectionID = connections[0]
	}

	return info, nil
}

// consensus Fetches the consensus state with the validator set.
func (n Node) consensus(ctx context.Context, client CosmosClient, height int64) (networktypes.Reward, error) {
	consensusState, err := client.ConsensusInfo(ctx, height)
	if err != nil {
		return networktypes.Reward{}, err
	}

	spnConsensusState := spntypes.NewConsensusState(
		consensusState.Timestamp,
		consensusState.NextValidatorsHash,
		consensusState.Root,
	)

	validators := make([]spntypes.Validator, len(consensusState.ValidatorSet.Validators))
	for i, validator := range consensusState.ValidatorSet.Validators {
		validators[i] = spntypes.NewValidator(
			base64.StdEncoding.EncodeToString(validator.PubKey.GetEd25519()),
			validator.ProposerPriority,
			validator.VotingPower,
		)
	}

	reward := networktypes.Reward{
		ConsensusState: spnConsensusState,
		ValidatorSet:   spntypes.NewValidatorSet(validators...),
		RevisionHeight: uint64(height),
	}

	return reward, nil
}

// connectionChannels fetches the chain connection channels by connection id.
func (n Node) connectionChannels(ctx context.Context, connectionID string) (channels []string, err error) {
	res, err := n.ibcChannelQuery.ConnectionChannels(ctx, &ibcchanneltypes.QueryConnectionChannelsRequest{
		Connection: connectionID,
	})
	if isNotFoundErr(err) {
		return channels, nil
	} else if err != nil {
		return nil, err
	}
	for _, channel := range res.Channels {
		channels = append(channels, channel.ChannelId)
	}
	return
}

// clientConnections fetches the chain client connections by client id.
func (n Node) clientConnections(ctx context.Context, clientID string) ([]string, error) {
	res, err := n.ibcConnQuery.ClientConnections(ctx, &ibcconntypes.QueryClientConnectionsRequest{
		ClientId: clientID,
	})
	if isNotFoundErr(err) {
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}
	return res.ConnectionPaths, err
}

// stakingParams fetches the staking module params.
func (n Node) stakingParams(ctx context.Context) (stakingtypes.Params, error) {
	res, err := n.stakingQuery.Params(ctx, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return stakingtypes.Params{}, err
	}
	return res.Params, nil
}

// consumerClientID fetches the consumer client id from the monitoring provider.
func (n Node) consumerClientID(ctx context.Context) (string, error) {
	res, err := n.monitoringProviderQuery.GetConsumerClientID(
		ctx, &monitoringptypes.QueryGetConsumerClientIDRequest{},
	)
	if isNotFoundErr(err) {
		return "", ErrObjectNotFound
	} else if err != nil {
		return "", err
	}
	return res.ConsumerClientId.ClientId, nil
}

// connectionChannelID fetches the consumer connection chnnael id from the monitoring provider.
func (n Node) connectionChannelID(ctx context.Context) (string, error) {
	res, err := n.monitoringProviderQuery.GetConnectionChannelID(
		ctx, &monitoringptypes.QueryGetConnectionChannelIDRequest{},
	)
	if isNotFoundErr(err) {
		return "", ErrObjectNotFound
	} else if err != nil {
		return "", err
	}
	return res.ConnectionChannelId.ChannelId, nil
}
