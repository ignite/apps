package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gookit/color"
	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/manifoldco/promptui"

	"github.com/ignite/apps/examples/spinner/cmd"
)

type app struct{}

func (app) Manifest(context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name:     "spinner",
		Commands: cmd.GetCommands(),
	}, nil
}

func (app) Execute(context.Context, *plugin.ExecutedCommand, plugin.ClientAPI) error {
	session := cliui.New(cliui.StartSpinnerWithText("Testing spinner..."))
	defer session.End()

	time.Sleep(time.Second * 5)
	session.StopSpinner()

	if err := session.AskConfirm("Testing question"); err != nil {
		if !errors.Is(err, promptui.ErrAbort) {
			return err
		}
	}

	session.StopSpinner()
	_ = session.Println(color.Blue.Sprintf("Asked question"))

	session.StartSpinner("Continue test spinner...")

	for i := time.Duration(0); i < 5; i++ {
		time.Sleep(time.Second * i)
		session.StartSpinner(fmt.Sprintf("Start again without stop %d", i))
	}

	for i := time.Duration(0); i < 5; i++ {
		time.Sleep(time.Second * 2 * i)
		session.StartSpinner(fmt.Sprintf("Start again with stop %d", i))
	}

	session.StopSpinner()
	_ = session.Printf(
		"%s %s\n",
		color.Green.Sprint("New test "),
		color.Yellow.Sprint("Color"),
	)
	session.StartSpinner("Continue test spinner...")

	time.Sleep(time.Second * 5)

	_ = session.Println(color.Green.Sprintf("Question without stop 1"))
	time.Sleep(time.Second * 5)
	_ = session.Println(color.Green.Sprintf("Question without stop 2"))
	time.Sleep(time.Second * 5)

	session.StopSpinner()

	_ = session.Println(color.Green.Sprintf("Question with stop 1"))

	session.StartSpinner("Finishing test spinner...")
	time.Sleep(time.Second * 5)
	session.StopSpinner()

	return nil
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
			"spinner": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
