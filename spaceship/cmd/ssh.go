package cmd

import (
	"context"

	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/spaceship/pkg/ssh"
)

// ExecuteSSH executes the ssh deploy subcommand.
func ExecuteSSH(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	var (
		host = "127.0.0.1"                        // arg host or URI
		user = "danilopantani"                    // flag user
		key  = "/Users/danilopantani/.ssh/id_rsa" // flag key
		// password = ""                          // flag password
		// port     = "22" // flag port
		// keyPassword = args[5] // flag key password
		// keyRaw      = args[6] // flag key raw
	)

	c, err := ssh.New(host, ssh.WithUser(user), ssh.WithKey(key))
	if err != nil {
		return err
	}
	if err := c.Connect(ctx); err != nil {
		return err
	}

	return nil
}
