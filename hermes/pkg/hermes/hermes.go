package hermes

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"time"

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
	FlagConfig           = "config"
	FlagJSON             = "json"
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

	// configs holds Generate configs.
	configs struct {
		flags  Flags
		config string
		stdOut io.Writer
		stdErr io.Writer
	}

	// Result represents the cli command result.
	Result struct {
		Result string `json:"result"`
		Status string `json:"status"`
	}

	// Log represents the cli command logs.
	Log struct {
		Timestamp time.Time `json:"timestamp"`
		Level     string    `json:"level"`
		Fields    Fields    `json:"fields"`
		ThreadId  string    `json:"threadId"`
	}

	// Fields represents the cli command result fields.
	Fields struct {
		Message string `json:"message"`
	}
)

// WithFlags assigns the command flags.
func WithFlags(flags Flags) Option {
	return func(c *configs) {
		c.flags = flags
	}
}

// WithConfigFile add a hermes config file.
func WithConfigFile(config string) Option {
	return func(c *configs) {
		c.config = config
	}
}

// WithStdOut add a std output.
func WithStdOut(stdOut io.Writer) Option {
	return func(c *configs) {
		c.stdOut = stdOut
	}
}

// WithStdErr add a std error output.
func WithStdErr(stdErr io.Writer) Option {
	return func(c *configs) {
		c.stdErr = stdErr
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

func (h *Hermes) AddKey(ctx context.Context, chainID, keyfile string, options ...Option) error {
	options = append(options, WithFlags(
		Flags{
			FlagChain:        chainID,
			FlagMnemonicFile: keyfile,
		},
	))
	return h.RunCmd(ctx, []string{string(cmdKeys)}, options...)
}

func (h *Hermes) AddMnemonic(ctx context.Context, chainID, mnemonic string, options ...Option) error {
	f, err := os.CreateTemp("", "hermes-key")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	if _, err := f.Write([]byte(mnemonic)); err != nil {
		return err
	}

	options = append(options, WithFlags(
		Flags{
			FlagChain:        chainID,
			FlagMnemonicFile: f.Name(),
		},
	))
	return h.RunCmd(ctx, []string{string(cmdKeys), string(cmdKeysAdd)}, options...)
}

func (h *Hermes) CreateClient(ctx context.Context, hostChain, referenceChain string, options ...Option) error {
	options = append(options, WithFlags(
		Flags{
			FlagHostChain:      hostChain,
			FlagReferenceChain: referenceChain,
		},
	))
	return h.RunCmd(ctx, []string{string(cmdCreate), string(cmdClient)}, options...)
}

func (h *Hermes) CreateConnection(ctx context.Context, chainA, clientA, clientB string, options ...Option) error {
	options = append(options, WithFlags(
		Flags{
			FlagChainA:  chainA,
			FlagClientA: clientA,
			FlagClientB: clientB,
		},
	))
	return h.RunCmd(ctx, []string{string(cmdCreate), string(cmdConnection)}, options...)
}

func (h *Hermes) CreateChannel(ctx context.Context, chainA, connA, portA, portB string, options ...Option) error {
	options = append(options, WithFlags(
		Flags{
			FlagChainA:      chainA,
			FlagConnectionA: connA,
			FlagPortA:       portA,
			FlagPortB:       portB,
		},
	))
	return h.RunCmd(ctx, []string{string(cmdCreate), string(cmdChannel)}, options...)
}

func (h *Hermes) QueryChannels(ctx context.Context, showCounterparty bool, chain string, options ...Option) error {
	flags := Flags{
		FlagChain: chain,
	}
	if showCounterparty {
		flags[FlagShowCounterparty] = true
	}
	options = append(options, WithFlags(flags))
	return h.RunCmd(ctx, []string{string(cmdQuery), string(cmdChannels)}, options...)
}

func (h *Hermes) Start(ctx context.Context, options ...Option) error {
	return h.RunCmd(ctx, []string{string(cmdStart)}, options...)
}

func (h *Hermes) RunCmd(ctx context.Context, args []string, options ...Option) error {
	c := configs{}
	for _, o := range options {
		o(&c)
	}

	cmd := args
	for flag, value := range c.flags {
		if v, ok := value.(bool); ok && v {
			cmd = append(cmd, fmt.Sprintf("--%s", flag))
		} else {
			cmd = append(cmd, fmt.Sprintf("--%s=%s", flag, value))
		}
	}

	stdOut := c.stdOut
	if stdOut == nil {
		stdOut = os.Stdout
	}

	stderr := c.stdErr
	if stderr == nil {
		stderr = os.Stderr
	}

	return h.Run(ctx, stdOut, stderr, c.config, cmd...)
}

func (h *Hermes) Run(ctx context.Context, stdOut, stdErr io.Writer, config string, args ...string) error {
	cmd := []string{h.path}
	cmd = append(cmd, fmt.Sprintf("--%s", FlagJSON))
	if config != "" {
		// the config flag should be added before the hermes subcommands
		cmd = append(cmd, fmt.Sprintf("--%s=%s", FlagConfig, config))
	}
	cmd = append(cmd, args...)
	return exec.Exec(ctx, cmd, exec.StepOption(step.Stdout(stdOut)), exec.StepOption(step.Stderr(stdErr)))
}
