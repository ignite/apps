package cmd

import (
	"time"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/bubbleconfirm"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/spaceship/pkg/ssh"
)

var ErrServerNotInitialized = errors.New("server not initialized")

const (
	flagPort         = "port"
	flagUser         = "user"
	flagUserPassword = "user-password"
	flagKey          = "key"
	flagRawKey       = "raw-key"
	flagKeyPassword  = "key-password"

	statusConnecting = "Connecting..."
)

func executeSSH(session *cliui.Session, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) (*ssh.SSH, error) {
	args := cmd.Args
	if len(args) < 1 {
		return nil, errors.New("must specify unless a uri host")
	}

	var (
		host = args[0]

		flags           = plugin.Flags(cmd.Flags)
		user, _         = flags.GetString(flagUser)
		userPassword, _ = flags.GetString(flagUserPassword)
		port, _         = flags.GetString(flagPort)
		key, _          = flags.GetString(flagKey)
		rawKey, _       = flags.GetString(flagRawKey)
		keyPassword, _  = flags.GetString(flagKeyPassword)
	)

	// Connect to the SSH.
	c, err := ssh.New(
		host,
		ssh.WithUser(user),
		ssh.WithPort(port),
		ssh.WithUserPassword(userPassword),
		ssh.WithKey(key),
		ssh.WithRawKey(rawKey),
		ssh.WithKeyPassword(keyPassword),
		ssh.WithWorkspace(chain.ChainId),
	)
	if err != nil {
		return nil, err
	}

	if c.NeedsUserPassword() {
		time.Sleep(10 * time.Millisecond) // Give some time to the spinner to start
		restart := session.PauseSpinner()
		if err := session.Ask(
			bubbleconfirm.NewQuestion(
				"Enter the SSH user password: ",
				&userPassword,
				bubbleconfirm.Required(),
				bubbleconfirm.HideAnswer(),
			),
		); err != nil {
			return nil, errors.Wrap(err, "you must provide a password")
		}
		restart()
	}

	return c, c.Connect()
}
