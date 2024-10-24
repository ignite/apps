package networkchain

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/cache"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosutil"
	cosmosgenesis "github.com/ignite/cli/v28/ignite/pkg/cosmosutil/genesis"
	"github.com/ignite/cli/v28/ignite/pkg/events"
	"github.com/ignite/cli/v28/ignite/pkg/jsonfile"
	launchtypes "github.com/ignite/network/x/launch/types"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/ignite/apps/network/network/networktypes"
)

// Prepare prepares the chain to be launched from genesis information.
func (c Chain) Prepare(
	ctx context.Context,
	cacheStorage cache.Storage,
	gi networktypes.GenesisInformation,
	rewardsInfo networktypes.Reward,
	spnChainID string,
	lastBlockHeight,
	consumerUnbondingTime int64,
) error {
	// chain initialization
	genesisPath, err := c.chain.GenesisPath()
	if err != nil {
		return err
	}

	_, err = os.Stat(genesisPath)

	switch {
	case os.IsNotExist(err):
		// if no config exists, perform a full initialization of the chain with a new validator key
		if err = c.Init(ctx, cacheStorage); err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		// if config and validator key already exists, build the chain and initialize the genesis
		if _, err := c.Build(ctx, cacheStorage); err != nil {
			return err
		}

		if err := c.initGenesis(ctx); err != nil {
			return err
		}
	}

	if err := c.buildGenesis(
		ctx,
		gi,
		rewardsInfo,
		spnChainID,
		lastBlockHeight,
		consumerUnbondingTime,
	); err != nil {
		return err
	}

	cmd, err := c.chain.Commands(ctx)
	if err != nil {
		return err
	}

	// ensure genesis has a valid format
	if err := cmd.ValidateGenesis(ctx); err != nil {
		return err
	}

	// reset the saved state in case the chain has been started before
	return cmd.UnsafeReset(ctx)
}

// buildGenesis builds the genesis for the chain from the launch approved requests.
func (c Chain) buildGenesis(
	ctx context.Context,
	gi networktypes.GenesisInformation,
	rewardsInfo networktypes.Reward,
	spnChainID string,
	lastBlockHeight,
	consumerUnbondingTime int64,
) error {
	c.ev.Send("Building the genesis", events.ProgressStart())

	addressPrefix, err := c.detectPrefix(ctx)
	if err != nil {
		return errors.Wrap(err, "error detecting chain prefix")
	}

	// apply genesis information to the genesis
	if err := c.applyGenesisAccounts(ctx, gi.GenesisAccounts, addressPrefix); err != nil {
		return errors.Wrap(err, "error applying genesis accounts to genesis")
	}
	if err := c.applyVestingAccounts(ctx, gi.VestingAccounts, addressPrefix); err != nil {
		return errors.Wrap(err, "error applying vesting accounts to genesis")
	}
	if err := c.applyGenesisValidators(ctx, gi.GenesisValidators); err != nil {
		return errors.Wrap(err, "error applying genesis validators to genesis")
	}

	genesisPath, err := c.chain.GenesisPath()
	if err != nil {
		return errors.Wrap(err, "genesis of the blockchain can't be read")
	}

	genesis, err := cosmosgenesis.FromPath(genesisPath)
	if err != nil {
		return errors.Wrap(err, "genesis of the blockchain can't be parsed")
	}

	// update chain ID and launch time
	if err := genesis.Update(
		jsonfile.WithKeyValue(cosmosgenesis.FieldChainID, c.id),
		jsonfile.WithKeyValueTimestamp(cosmosgenesis.FieldGenesisTime, c.launchTime.Unix()),
	); err != nil {
		return errors.Wrap(err, "genesis cannot be updated")
	}

	// update reward related fields if the testnet is incentivized (with a last block height for reward distribution)
	if lastBlockHeight > 0 {
		if err := genesis.Update(
			jsonfile.WithKeyValue(cosmosgenesis.FieldConsumerChainID, spnChainID),
			jsonfile.WithKeyValueInt(cosmosgenesis.FieldLastBlockHeight, lastBlockHeight),
			jsonfile.WithKeyValue(cosmosgenesis.FieldConsensusTimestamp, rewardsInfo.ConsensusState.Timestamp),
			jsonfile.WithKeyValue(cosmosgenesis.FieldConsensusNextValidatorsHash, rewardsInfo.ConsensusState.NextValidatorsHash),
			jsonfile.WithKeyValue(cosmosgenesis.FieldConsensusRootHash, rewardsInfo.ConsensusState.Root.Hash),
			jsonfile.WithKeyValueInt(cosmosgenesis.FieldConsumerUnbondingPeriod, consumerUnbondingTime),
			jsonfile.WithKeyValueUint(cosmosgenesis.FieldConsumerRevisionHeight, rewardsInfo.RevisionHeight),
		); err != nil {
			return errors.Wrap(err, "genesis cannot be updated for reward related fields")
		}
	}

	if err := applyParamChanges(genesis, gi.ParamChanges); err != nil {
		return fmt.Errorf("error applying param changes to genesis: %w", err)
	}

	c.ev.Send("Genesis built", events.ProgressFinish())

	return nil
}

// applyGenesisAccounts adds the genesis account into the genesis using the chain CLI.
func (c Chain) applyGenesisAccounts(
	ctx context.Context,
	genesisAccs []networktypes.GenesisAccount,
	addressPrefix string,
) error {
	var err error

	cmd, err := c.chain.Commands(ctx)
	if err != nil {
		return err
	}

	for _, acc := range genesisAccs {
		// change the address prefix to the target chain prefix
		acc.Address, err = cosmosutil.ChangeAddressPrefix(acc.Address, addressPrefix)
		if err != nil {
			return err
		}

		// call the add genesis account CLI command
		err = cmd.AddGenesisAccount(ctx, acc.Address, acc.Coins.String())
		if err != nil {
			return err
		}
	}

	return nil
}

// applyVestingAccounts adds the genesis vesting account into the genesis using the chain CLI.
func (c Chain) applyVestingAccounts(
	ctx context.Context,
	vestingAccs []networktypes.VestingAccount,
	addressPrefix string,
) error {
	cmd, err := c.chain.Commands(ctx)
	if err != nil {
		return err
	}

	for _, acc := range vestingAccs {
		acc.Address, err = cosmosutil.ChangeAddressPrefix(acc.Address, addressPrefix)
		if err != nil {
			return err
		}

		// call the add genesis account CLI command with delayed vesting option
		err = cmd.AddVestingAccount(
			ctx,
			acc.Address,
			acc.TotalBalance.String(),
			acc.Vesting.String(),
			acc.EndTime,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// applyGenesisValidators gathers the validator gentxs into the genesis and adds peers in config.
func (c Chain) applyGenesisValidators(ctx context.Context, genesisVals []networktypes.GenesisValidator) error {
	// no validator
	if len(genesisVals) == 0 {
		return nil
	}

	// reset the gentx directory
	gentxDir, err := c.chain.GentxsPath()
	if err != nil {
		return err
	}
	if err := os.RemoveAll(gentxDir); err != nil {
		return err
	}
	if err := os.MkdirAll(gentxDir, 0o700); err != nil {
		return err
	}

	// write gentxs
	for i, val := range genesisVals {
		gentxPath := filepath.Join(gentxDir, fmt.Sprintf("gentx%d.json", i))
		if err = os.WriteFile(gentxPath, val.Gentx, 0o666); err != nil {
			return err
		}
	}

	// gather gentxs
	cmd, err := c.chain.Commands(ctx)
	if err != nil {
		return err
	}
	if err := cmd.CollectGentxs(ctx); err != nil {
		return err
	}

	return c.updateConfigFromGenesisValidators(genesisVals)
}

// applyParamChanges applies the param changes into the genesis.
func applyParamChanges(
	genesis *cosmosgenesis.Genesis,
	paramChanges []networktypes.ParamChange,
) error {
	changes := make([]jsonfile.UpdateFileOption, len(paramChanges))

	for i, pc := range paramChanges {
		changes[i] = jsonfile.WithKeyValueByte(cosmosgenesis.ModuleParamField(pc.Module, pc.Param), pc.Value)
	}

	if err := genesis.Update(changes...); err != nil {
		return errors.Wrap(err, "failed to apply param change to genesis")
	}

	return nil
}

// updateConfigFromGenesisValidators adds the peer addresses into the config.toml of the chain.
func (c Chain) updateConfigFromGenesisValidators(genesisVals []networktypes.GenesisValidator) error {
	var (
		p2pAddresses    []string
		tunnelAddresses []TunneledPeer
	)
	for i, val := range genesisVals {
		if !networktypes.VerifyPeerFormat(val.Peer) {
			return errors.Errorf("invalid peer: %s", val.Peer.Id)
		}
		switch conn := val.Peer.Connection.(type) {
		case *launchtypes.Peer_TcpAddress:
			p2pAddresses = append(p2pAddresses, fmt.Sprintf("%s@%s", val.Peer.Id, conn.TcpAddress))
		case *launchtypes.Peer_HttpTunnel:
			tunneledPeer := TunneledPeer{
				Name:      conn.HttpTunnel.Name,
				Address:   conn.HttpTunnel.Address,
				NodeID:    val.Peer.Id,
				LocalPort: strconv.Itoa(i + 22000),
			}
			tunnelAddresses = append(tunnelAddresses, tunneledPeer)
			p2pAddresses = append(p2pAddresses, fmt.Sprintf("%s@127.0.0.1:%s", tunneledPeer.NodeID, tunneledPeer.LocalPort))
		default:
			return fmt.Errorf("invalid peer type")
		}
	}

	if len(p2pAddresses) > 0 {
		// set persistent peers
		configPath, err := c.chain.ConfigTOMLPath()
		if err != nil {
			return err
		}
		configToml, err := toml.LoadFile(configPath)
		if err != nil {
			return err
		}
		configToml.Set("p2p.persistent_peers", strings.Join(p2pAddresses, ","))
		if err != nil {
			return err
		}

		// if there are tunneled peers they will be connected with tunnel clients via localhost,
		// so we need to allow to have few nodes with the same ip
		if len(tunnelAddresses) > 0 {
			configToml.Set("p2p.allow_duplicate_ip", true)
		}

		// save config.toml file
		configTomlFile, err := os.OpenFile(configPath, os.O_RDWR|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
		defer configTomlFile.Close()

		if _, err = configToml.WriteTo(configTomlFile); err != nil {
			return err
		}
	}

	if len(tunnelAddresses) > 0 {
		tunneledPeersConfigPath, err := c.SPNConfigPath()
		if err != nil {
			return err
		}

		if err = SetSPNConfig(Config{
			TunneledPeers: tunnelAddresses,
		}, tunneledPeersConfigPath); err != nil {
			return err
		}
	}
	return nil
}
