package networkchain

import (
	"context"
	"fmt"

	cosmosgenesis "github.com/ignite/cli/v28/ignite/pkg/cosmosutil/genesis"
	"github.com/ignite/cli/v28/ignite/pkg/events"
	"github.com/pkg/errors"

	"github.com/ignite/apps/network/network/networktypes"
)

// CheckRequestChangeParam builds the genesis for the chain from the launch approved requests.
func (c Chain) CheckRequestChangeParam(
	ctx context.Context,
	module,
	param string,
	value []byte,
) error {
	c.ev.Send("Checking the param change", events.ProgressStart())

	if err := c.initGenesis(ctx); err != nil {
		return err
	}

	genesisPath, err := c.chain.GenesisPath()
	if err != nil {
		return errors.Wrap(err, "genesis of the blockchain can't be read")
	}

	genesis, err := cosmosgenesis.FromPath(genesisPath)
	if err != nil {
		return errors.Wrap(err, "genesis of the blockchain can't be parsed")
	}

	pc := []networktypes.ParamChange{
		{
			Module: module,
			Param:  param,
			Value:  value,
		},
	}

	if err := applyParamChanges(genesis, pc); err != nil {
		return fmt.Errorf("error applying param changes to genesis: %w", err)
	}

	cmd, err := c.chain.Commands(ctx)
	if err != nil {
		return err
	}

	// ensure genesis has a valid format
	if err := cmd.ValidateGenesis(ctx); err != nil {
		return fmt.Errorf("invalid parameter change requested: %w", err)
	}

	c.ev.Send("Param change verified", events.ProgressFinish())

	return nil
}
