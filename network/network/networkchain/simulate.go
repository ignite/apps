package networkchain

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/v28/ignite/pkg/availableport"
	"github.com/ignite/cli/v28/ignite/pkg/cache"
	"github.com/ignite/cli/v28/ignite/pkg/events"
	"github.com/ignite/cli/v28/ignite/pkg/httpstatuschecker"
	"github.com/ignite/cli/v28/ignite/pkg/xurl"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/ignite/apps/network/network/networktypes"
)

// listeningTimeout is the maximum time to wait for the chain start.
const listeningTimeout = time.Minute * 1

// validatorSetNilErrorMessages contains variation for SDK validator set empty error message/
// to detect if the chain fails to start because the validator set is nil.
var validatorSetNilErrorMessages = []string{
	"validator set is nil in genesis and still empty after InitChain",
	"validator set is empty after InitGenesis",
}

// SimulateRequests simulates the genesis creation and the start of the network from the provided requests.
func (c Chain) SimulateRequests(
	ctx context.Context,
	cacheStorage cache.Storage,
	gi networktypes.GenesisInformation,
	reqs []networktypes.Request,
) (err error) {
	c.ev.Send("Verifying requests format", events.ProgressStart())
	for _, req := range reqs {
		// static verification of the request
		if err := networktypes.VerifyRequest(req); err != nil {
			return err
		}

		// apply the request to the genesis information
		gi, err = gi.ApplyRequest(req)
		if err != nil {
			return err
		}
	}
	c.ev.Send("Requests format verified", events.ProgressFinish())

	// prepare the chain with the requests
	if err := c.Prepare(
		ctx,
		cacheStorage,
		gi,
		networktypes.Reward{RevisionHeight: 1},
		networktypes.SPNChainID,
		1,
		2,
	); err != nil {
		return err
	}

	c.ev.Send("Trying starting the network with the requests", events.ProgressStart())
	if err := c.simulateChainStart(ctx); err != nil {
		return err
	}
	c.ev.Send("The network can be started", events.ProgressFinish())

	return nil
}

// simulateChainStart simulates and verify the chain start by starting it with a simulation config
// and checking if the gentxs execution is successful.
func (c Chain) simulateChainStart(ctx context.Context) error {
	cmd, err := c.chain.Commands(ctx)
	if err != nil {
		return err
	}

	// set the config with random ports to test the start command
	rpcAddr, err := c.setSimulationConfig()
	if err != nil {
		return err
	}

	// verify that the chain can be started with a valid genesis
	ctx, cancel := context.WithTimeout(ctx, listeningTimeout)
	exit := make(chan error)

	// routine to check the app is listening
	go func() {
		defer cancel()
		exit <- isChainListening(ctx, rpcAddr)
	}()

	// check an error is realted to validator set empty
	checkValidatorSetEmptyError := func(err error) bool {
		if err != nil {
			for _, errStr := range validatorSetNilErrorMessages {
				if strings.Contains(err.Error(), errStr) {
					return true
				}
			}
		}

		return false
	}

	// routine chain start
	go func() {
		// if the error is validator set is nil, it means the genesis didn't get broken after an applied request
		// the genesis was correctly generated but there is no gentxs so far
		// so we don't consider it as an error making requests to verify as invalid
		err := cmd.Start(ctx)
		if checkValidatorSetEmptyError(err) {
			err = nil
		}
		exit <- errors.Wrap(err, "the chain failed to start")
	}()

	return <-exit
}

// setSimulationConfig sets in the config random available ports to allow check if the chain network can start.
func (c Chain) setSimulationConfig() (string, error) {
	// generate random server ports and servers list
	ports, err := availableport.Find(5)
	if err != nil {
		return "", err
	}
	genAddr := func(port uint) string {
		return fmt.Sprintf("0.0.0.0:%d", port)
	}

	// updating app toml
	appPath, err := c.AppTOMLPath()
	if err != nil {
		return "", err
	}
	config, err := toml.LoadFile(appPath)
	if err != nil {
		return "", err
	}

	apiAddr, err := xurl.TCP(genAddr(ports[0]))
	if err != nil {
		return "", err
	}

	config.Set("api.enable", true)
	config.Set("api.enabled-unsafe-cors", true)
	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("api.address", apiAddr)
	config.Set("grpc.address", genAddr(ports[1]))
	gas := sdktypes.NewInt64Coin(networktypes.SPNDenom, 0)
	config.Set("minimum-gas-prices", gas.String())

	file, err := os.OpenFile(appPath, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := config.WriteTo(file); err != nil {
		return "", err
	}

	// updating config toml
	configPath, err := c.ConfigTOMLPath()
	if err != nil {
		return "", err
	}
	config, err = toml.LoadFile(configPath)
	if err != nil {
		return "", err
	}

	rpcAddr, err := xurl.TCP(genAddr(ports[2]))
	if err != nil {
		return "", err
	}

	p2pAddr, err := xurl.TCP(genAddr(ports[3]))
	if err != nil {
		return "", err
	}

	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("consensus.timeout_commit", "1s")
	config.Set("consensus.timeout_propose", "1s")
	config.Set("rpc.laddr", rpcAddr)
	config.Set("p2p.laddr", p2pAddr)
	config.Set("rpc.pprof_laddr", genAddr(ports[4]))

	file, err = os.OpenFile(configPath, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = config.WriteTo(file)

	return genAddr(ports[2]), err
}

// isChainListening checks if the chain is listening for RPC queries on the specified address.
func isChainListening(ctx context.Context, rpcAddr string) error {
	checkAlive := func() error {
		addr, err := xurl.HTTP(rpcAddr)
		if err != nil {
			return fmt.Errorf("invalid rpc address format %s: %w", rpcAddr, err)
		}

		ok, err := httpstatuschecker.Check(ctx, fmt.Sprintf("%s/health", addr))
		if err == nil && !ok {
			err = errors.New("app is not online")
		}
		return err
	}
	return backoff.Retry(checkAlive, backoff.WithContext(backoff.NewConstantBackOff(time.Second), ctx))
}
