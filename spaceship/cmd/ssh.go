package cmd

import (
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/spaceship/pkg/ssh"
)

var ErrServerNotInitialized = errors.New("server not initialized")

const (
	flagPort        = "port"
	flagUser        = "user"
	flagPassword    = "password"
	flagKey         = "key"
	flagRawKey      = "raw-key"
	flagKeyPassword = "key-password"
	flagInitChain   = "init-chain"
	flagFaucet      = "faucet"
	flagFaucetPort  = "faucet-port"
	flagLines       = "lines"
	flagRealTime    = "real-time"

	statusConnecting = "Connecting..."
)

func executeSSH(cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) (*ssh.SSH, error) {
	args := cmd.Args
	if len(args) < 1 {
		return nil, errors.New("must specify unless a uri host")
	}

	var (
		host = args[0]

		flags          = plugin.Flags(cmd.Flags)
		user, _        = flags.GetString(flagUser)
		port, _        = flags.GetString(flagPort)
		password, _    = flags.GetString(flagPassword)
		key, _         = flags.GetString(flagKey)
		rawKey, _      = flags.GetString(flagRawKey)
		keyPassword, _ = flags.GetString(flagKeyPassword)
	)

	// Connect to the SSH.
	c, err := ssh.New(
		host,
		ssh.WithUser(user),
		ssh.WithPort(port),
		ssh.WithPassword(password),
		ssh.WithKey(key),
		ssh.WithRawKey(rawKey),
		ssh.WithKeyPassword(keyPassword),
		ssh.WithWorkspace(chain.ChainId),
	)
	if err != nil {
		return nil, err
	}

	return c, c.Connect()
}
