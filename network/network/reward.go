package network

import (
	"context"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v28/ignite/pkg/events"
	rewardtypes "github.com/ignite/network/x/reward/types"

	"github.com/ignite/apps/network/network/networktypes"
)

// SetReward set a chain reward.
func (n Network) SetReward(ctx context.Context, launchID uint64, lastRewardHeight int64, coins sdk.Coins) error {
	n.ev.Send(
		fmt.Sprintf("Setting reward %s to the chain %d at height %d", coins, launchID, lastRewardHeight),
		events.ProgressStart(),
	)

	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return err
	}

	msg := rewardtypes.NewMsgSetRewards(
		addr,
		launchID,
		coins,
		lastRewardHeight,
	)
	res, err := n.cosmos.BroadcastTx(ctx, n.account, msg)
	if err != nil {
		return err
	}

	var setRewardRes rewardtypes.MsgSetRewardsResponse
	if err := res.Decode(&setRewardRes); err != nil {
		return err
	}

	if setRewardRes.PreviousCoins.Empty() {
		n.ev.Send("The reward pool was empty", events.Icon(icons.Info), events.ProgressFinish())
	} else {
		n.ev.Send(
			fmt.Sprintf("Previous reward pool %s at height %d is overwritten", coins, lastRewardHeight),
			events.Icon(icons.Info),
			events.ProgressFinish(),
		)
	}

	if setRewardRes.NewCoins.Empty() {
		n.ev.Send("The reward pool is removed", events.ProgressFinish())
	} else {
		n.ev.Send(
			fmt.Sprintf(
				"%s will be distributed to validators at height %d. The chain %d is now an incentivized testnet",
				coins,
				lastRewardHeight,
				launchID,
			),
			events.ProgressFinish(),
		)
	}
	return nil
}

// RewardsInfo Fetches the consensus state with the validator set,
// the unbounding time, and the last block height from chain rewards.
func (n Network) RewardsInfo(
	ctx context.Context,
	launchID uint64,
	height int64,
) (
	rewardsInfo networktypes.Reward,
	lastRewardHeight int64,
	unboundingTime int64,
	err error,
) {
	rewardsInfo, err = n.node.consensus(ctx, n.cosmos, height)
	if err != nil {
		return rewardsInfo, 0, 0, err
	}

	stakingParams, err := n.node.stakingParams(ctx)
	if err != nil {
		return rewardsInfo, 0, 0, err
	}
	unboundingTime = int64(stakingParams.UnbondingTime.Seconds())

	chainReward, err := n.ChainReward(ctx, launchID)
	if errors.Is(err, ErrObjectNotFound) {
		return rewardsInfo, 1, unboundingTime, nil
	} else if err != nil {
		return rewardsInfo, 0, 0, err
	}
	lastRewardHeight = chainReward.LastRewardHeight

	return
}
