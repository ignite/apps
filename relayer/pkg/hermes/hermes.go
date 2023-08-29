package hermes

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/localfs"
	"github.com/ignite/ignite-files/hermes"
)

const (
	FlagHostChain        = "host-chain"
	FlagReferenceChain   = "reference-chain"
	FlagChainA           = "a-chain"
	FlagChainB           = "b-chain"
	FlagClientA          = "a-client"
	FlagClientB          = "b-client"
	FlagConnectionA      = "a-connection"
	FlagConnectionB      = "b-connection"
	FlagPortA            = "a-port"
	FlagPortB            = "b-port"
	FlagShowCounterparty = "show-counterparty"
	FlagChain            = "chain"
	FlagMnemonicFile     = "mnemonic-file"
)

const (
	// CommandCreate is the Hermes create command.
	cmdCreate cmdName = "create"

	// CommandQuery is the Hermes query command.
	cmdQuery cmdName = "query"

	// CommandKeys is the Hermes keys command.
	cmdKeys cmdName = "keys"

	// CommandStart is the Hermes start command.
	cmdStart cmdName = "start"

	// CommandClient is the Hermes create client command.
	cmdClient subCmd = "client"

	// CommandConnection is the Hermes create connection command.
	cmdConnection subCmd = "connection"

	// CommandChannel is the Hermes create channel command.
	cmdChannel subCmd = "channel"

	// CommandChannels  is the Hermes query channels command.
	cmdChannels subCmd = "channels"

	// CommandKeysAdd is the Hermes keys add command.
	cmdKeysAdd subCmd = "add"
)

type (
	// Flags represents the Hermes run flags.
	Flags map[string]interface{}
	// cmdName represents a high level command under Hermes.
	cmdName string
	// SubCommand represents the sub command under Hermes.
	subCmd string

	Hermes struct {
		path    string
		binary  []byte
		cleanup func()
	}

	// Option configures Generate configs.
	Option func(*configs)

	// Configs holds Generate configs.
	configs struct {
		flags Flags
	}
)

// WithFlags assigns the command flags.
func WithFlags(flags Flags) Option {
	return func(c *configs) {
		c.flags = flags
	}
}

// New returns the hermes binary executable.
func New() (*Hermes, error) {
	// untar the binary.
	gzr, err := gzip.NewReader(bytes.NewReader(hermes.Binary()))
	if err != nil {
		panic(err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	if _, err := tr.Next(); err != nil {
		return nil, err
	}

	binary, err := io.ReadAll(tr)
	if err != nil {
		return nil, err
	}

	path, cleanup, err := localfs.SaveBytesTemp(binary, "hermes", 0o755)
	if err != nil {
		return nil, err
	}

	return &Hermes{
		path:    path,
		binary:  binary,
		cleanup: cleanup,
	}, nil
}

func (h *Hermes) Cleanup() error {
	h.cleanup()
	h.binary = nil
	return os.RemoveAll(h.path)
}

func (h *Hermes) AddKey(ctx context.Context, chainID, keyfile string) ([]byte, error) {
	return h.RunCmd(ctx, cmdKeys, "", WithFlags(
		Flags{
			FlagChain:        chainID,
			FlagMnemonicFile: keyfile,
		},
	))
}

func (h *Hermes) AddMnemonic(ctx context.Context, chainID, mnemonic string) ([]byte, error) {
	f, err := os.CreateTemp("", "hermes-key")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())
	if _, err := f.Write([]byte(mnemonic)); err != nil {
		return nil, err
	}
	return h.RunCmd(ctx, cmdKeys, cmdKeysAdd, WithFlags(
		Flags{
			FlagChain:        chainID,
			FlagMnemonicFile: f.Name(),
		},
	))
}

func (h *Hermes) CreateClient(ctx context.Context, hostChain, referenceChain string) ([]byte, error) {
	return h.RunCmd(ctx, cmdCreate, cmdClient, WithFlags(
		Flags{
			FlagHostChain:      hostChain,
			FlagReferenceChain: referenceChain,
		},
	))
}

func (h *Hermes) CreateConnection(ctx context.Context, chainA, clientA, clientB string) ([]byte, error) {
	return h.RunCmd(ctx, cmdCreate, cmdConnection, WithFlags(
		Flags{
			FlagChainA:  chainA,
			FlagClientA: clientA,
			FlagClientB: clientB,
		}))
}

func (h *Hermes) CreateChannel(ctx context.Context, chainA, connA, portA, portB string) ([]byte, error) {
	return h.RunCmd(ctx, cmdCreate, cmdChannel, WithFlags(
		Flags{
			FlagChainA:      chainA,
			FlagConnectionA: connA,
			FlagPortA:       portA,
			FlagPortB:       portB,
		},
	))
}

func (h *Hermes) QueryChannels(ctx context.Context, showCounterparty bool, chain string) ([]byte, error) {
	flags := Flags{
		FlagChain: chain,
	}
	if showCounterparty {
		flags[FlagShowCounterparty] = true
	}
	return h.RunCmd(ctx, cmdQuery, cmdChannels, WithFlags(flags))
}

func (h *Hermes) Start(ctx context.Context, options ...Option) ([]byte, error) {
	return h.RunCmd(ctx, cmdStart, "", options...)
}

func (h *Hermes) RunCmd(ctx context.Context, command cmdName, subCommand subCmd, options ...Option) ([]byte, error) {
	c := configs{}
	for _, o := range options {
		o(&c)
	}

	cmd := []string{h.path, string(command)}
	if subCommand != "" {
		cmd = append(cmd, string(subCommand))
	}
	for flag, value := range c.flags {
		if v, ok := value.(bool); ok && v {
			cmd = append(cmd, fmt.Sprintf("--%s", flag))
		} else {
			cmd = append(cmd, fmt.Sprintf("--%s=%s", flag, value))
		}
	}

	// execute the command.
	buf := bytes.Buffer{}
	return buf.Bytes(), exec.Exec(ctx, cmd, exec.StepOption(step.Stdout(&buf)))
}
