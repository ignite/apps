package main

import (
	"context"
	"io"
	"os"

	"github.com/creack/pty"
	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/hermes/cmd"
)

type app struct{}

func (app) Manifest(context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name:     "hermes",
		Commands: cmd.GetCommands(),
	}, nil
}

// stdOutRedirect manages the redirection of stdout/stderr with a PTY.
func stdOutRedirect() error {
	// Save the current stdOut.
	out := os.Stdout

	// Create a new pseudo-terminal.
	pty, tty, err := pty.Open()
	if err != nil {
		return errors.Errorf("Failed to create pty/tty: %v", err)
	}

	// Redirect os.Stdout and os.Stderr to the TTY.
	os.Stdout = tty

	// Start a goroutine to forward the output from the PTY to the real terminal.
	go func() {
		_, _ = io.Copy(out, pty) // Copy all output from the PTY to the real terminal.
	}()
	return nil
}

func (app) Execute(ctx context.Context, c *plugin.ExecutedCommand, _ plugin.ClientAPI) error {
	// The CLIUI uses the github.com/briandowns/spinner library, which checks if the output is a terminal since it can only display in a terminal.
	// The Ignite app adds a stdout layer, causing a misunderstanding about whether the output is a terminal.
	// To prevent the spinner from being printed in the plugin, this is a workaround for the library.
	// A better solution would involve either changing the spinner logic or improving the detection mechanism.
	// Reference: https://github.com/briandowns/spinner/blob/master/spinner.go#L502-L506
	if err := stdOutRedirect(); err != nil {
		return err
	}

	// Remove the three two elements "ignite", "relayer" and "hermes" from OsArgs.
	args := c.OsArgs[3:]
	switch args[0] {
	case "configure":
		return cmd.ConfigureHandler(ctx, c)
	case "exec":
		return cmd.ExecuteHandler(ctx, c)
	case "start":
		return cmd.StartHandler(ctx, c)
	case "keys":
		switch args[1] {
		case "add":
			return cmd.KeysAddMnemonicHandler(ctx, c)
		case "file":
			return cmd.KeysAddFileHandler(ctx, c)
		case "list":
			return cmd.KeysListHandler(ctx, c)
		case "delete":
			return cmd.KeysDeleteHandler(ctx, c)
		default:
			return errors.Errorf("unknown keys command: %s", args[1])
		}
	default:
		return errors.Errorf("unknown command: %s", args[0])
	}
}

func (app) ExecuteHookPre(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookPost(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookCleanUp(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
	return nil
}

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins: map[string]hplugin.Plugin{
			"hermes": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
