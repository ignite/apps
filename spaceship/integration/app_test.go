package integration_test

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	pluginsconfig "github.com/ignite/cli/v29/ignite/config/plugins"
	"github.com/ignite/cli/v29/ignite/pkg/availableport"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/plugin"
	envtest "github.com/ignite/cli/v29/integration"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

const privateKey = `# random generated pk
-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----`

func newServer(t *testing.T, port uint, user, password string) {
	t.Helper()

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == user && string(pass) == password {
				return nil, nil
			}
			return nil, errors.Errorf("password rejected for %q", c.User())
		},
	}

	private, err := ssh.ParsePrivateKey([]byte(privateKey))
	require.NoError(t, err, "failed to parse private key")

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	require.NoError(t, err, "failed to listen for connection")

	go func() {
		nConn, err := listener.Accept()
		require.NoError(t, err, "failed to accept incoming connection")

		// Before use, a handshake must be performed on the incoming
		// net.Conn.
		_, chans, reqs, err := ssh.NewServerConn(nConn, config)
		require.NoError(t, err, "failed to handshake")

		// The incoming Request channel must be serviced.
		go ssh.DiscardRequests(reqs)

		// Service the incoming channel.
		for newChannel := range chans {
			if newChannel.ChannelType() != "session" {
				newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
				continue
			}

			channel, requests, err := newChannel.Accept()
			require.NoError(t, err, "could not accept channel")

			go func(in <-chan *ssh.Request) {
				for req := range in {
					switch req.Type {
					case "exec":
						// just return error 0 without exec.
						channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
					}
					req.Reply(req.Type == "exec", nil)
				}
			}(requests)

			term := terminal.NewTerminal(channel, "> ")
			go func() {
				defer channel.Close()
				for {
					line, err := term.ReadLine()
					if err != nil {
						break
					}
					fmt.Println(line)
				}
			}()
		}
	}()
}

func TestSpaceship(t *testing.T) {
	t.Skip("find a better way to deploy a SSH server for tests")

	var (
		require  = require.New(t)
		env      = envtest.New(t)
		app      = env.Scaffold("spaceship-app")
		user     = "spaceship"
		password = "spaceship_pass"
	)

	port, err := availableport.Find(1)
	require.NoError(err)

	newServer(t, port[0], user, password)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "spaceship")

	env.Must(env.Exec("install spaceship app locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", pluginPath),
			step.Workdir(app.SourcePath()),
		)),
	))

	// One local plugin expected
	assertLocalPlugins(t, app, []pluginsconfig.Plugin{{Path: pluginPath}})
	assertGlobalPlugins(t, nil)

	env.Must(env.Exec("run spaceship",
		step.NewSteps(step.New(
			step.Workdir(app.SourcePath()),
			step.Stdout(os.Stdout),
			step.Stderr(os.Stderr),
			step.Stdin(os.Stdin),
			step.Exec(
				envtest.IgniteApp,
				"spaceship",
				"deploy",
				"127.0.0.1",
				"-u", user,
				"--password", password,
				"-p", strconv.Itoa(int(port[0])),
			),
		)),
	))
}

func assertLocalPlugins(t *testing.T, app envtest.App, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfg, err := pluginsconfig.ParseDir(app.SourcePath())
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected local apps")
}

func assertGlobalPlugins(t *testing.T, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfgPath, err := plugin.PluginsPath()
	require.NoError(t, err)
	cfg, err := pluginsconfig.ParseDir(cfgPath)
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected global apps")
}
