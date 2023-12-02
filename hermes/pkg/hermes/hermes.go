package hermes

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
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
	FlagKeyName          = "key-name"
	FlagConfig           = "config"
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

	// cmdKeysList is the Hermes keys list command.
	cmdKeysList subCmd = "list"

	// cmdKeysDelete is the Hermes keys delete command.
	cmdKeysDelete subCmd = "delete"

	// ResultSuccess is the api result status success.
	ResultSuccess = "success"

	// ResultError is the api result status error.
	ResultError = "error"
)

// ErrResult indicates that Hermes binary returned an error.
var ErrResult = errors.New("result error")

type (
	// Flags represents the Hermes run flags.
	Flags map[string]interface{}
	// cmdName represents a high level command under Hermes.
	cmdName string
	// SubCommand represents the sub command under Hermes.
	subCmd string

	// Hermes represents the hermes binary structure.
	Hermes struct {
		path    string
		cleanup func()
	}

	// Option configures Generate configs.
	Option func(*configs)

	// configs holds Generate configs.
	configs struct {
		flags  Flags
		config string
		args   []string
		stdout io.Writer
		stderr io.Writer
	}
)

type (
	// Result represents the command result.
	Result struct {
		Result json.RawMessage `json:"result"`
		Status string          `json:"status"`
	}

	// KeysListResult represents the result of the keys list command.
	KeysListResult struct {
		Wallet Wallet `json:"wallet"`
	}

	// Wallet represents the wallet from a hermes key.
	Wallet struct {
		Account     string `json:"account"`
		Address     []byte `json:"address"`
		AddressType string `json:"address_type"`
		PrivateKey  string `json:"private_key"`
		PublicKey   string `json:"public_key"`
	}

	// ClientResult represents the result of the create client command.
	ClientResult struct {
		CreateClient CreateClient `json:"CreateClient"`
	}

	// CreateClient represents the result of the create client command.
	CreateClient struct {
		ClientID        string          `json:"client_id"`
		ClientType      string          `json:"client_type"`
		ConsensusHeight ConsensusHeight `json:"consensus_height"`
	}

	// ConsensusHeight represents the result consensus height.
	ConsensusHeight struct {
		RevisionHeight int `json:"revision_height"`
		RevisionNumber int `json:"revision_number"`
	}

	// ConnectionResult represents the result of the create connection command.
	ConnectionResult struct {
		ASide           Side   `json:"a_side"`
		BSide           Side   `json:"b_side"`
		ConnectionDelay Time   `json:"connection_delay"`
		DelayPeriod     Time   `json:"delay_period"`
		Ordering        string `json:"ordering"`
	}

	// Side represents the connection side.
	Side struct {
		ChannelID    string      `json:"channel_id"`
		ClientID     string      `json:"client_id"`
		ConnectionID string      `json:"connection_id"`
		PortID       string      `json:"port_id"`
		Version      interface{} `json:"version"`
	}

	// Time represents the time.
	Time struct {
		Nanos int `json:"nanos"`
		Secs  int `json:"secs"`
	}

	// ChannelResult represents the result of the create channel command.
	ChannelResult struct {
		ChainIDA string `json:"chain_id_a"`
		ChainIDB string `json:"chain_id_b"`
		ChannelA string `json:"channel_a"`
		ChannelB string `json:"channel_b"`
		PortA    string `json:"port_a"`
		PortB    string `json:"port_b"`
	}
)

// WithArgs assigns the command args.
func WithArgs(args ...string) Option {
	return func(c *configs) {
		c.args = args
	}
}

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
		c.stdout = stdOut
	}
}

// WithStdErr add a std error output.
func WithStdErr(stdErr io.Writer) Option {
	return func(c *configs) {
		c.stderr = stdErr
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
		cleanup: cleanup,
	}, nil
}

// Cleanup clean the temporary Hermes binary.
func (h *Hermes) Cleanup() error {
	h.cleanup()
	return os.RemoveAll(h.path)
}

// AddKey adds a new key file into the Hermes.
func (h *Hermes) AddKey(ctx context.Context, chainID, keyfile string, options ...Option) error {
	options = append(
		options,
		WithFlags(Flags{
			FlagChain:        chainID,
			FlagMnemonicFile: keyfile,
		}),
		WithArgs(string(cmdKeys), string(cmdKeysAdd)),
	)
	return h.Run(ctx, options...)
}

// AddMnemonic creates a new temporary key file based on the mnemonic and add into the Hermes.
func (h *Hermes) AddMnemonic(ctx context.Context, chainID, mnemonic string, options ...Option) error {
	f, err := os.CreateTemp("", "hermes-key")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	if _, err := f.Write([]byte(mnemonic)); err != nil {
		return err
	}

	options = append(
		options,
		WithFlags(Flags{
			FlagChain:        chainID,
			FlagMnemonicFile: f.Name(),
		}),
		WithArgs(string(cmdKeys), string(cmdKeysAdd)),
	)
	return h.Run(ctx, options...)
}

// KeysList list all available Hermes keys.
func (h *Hermes) KeysList(ctx context.Context, chainID string, options ...Option) error {
	options = append(
		options,
		WithFlags(Flags{FlagChain: chainID}),
		WithArgs(string(cmdKeys), string(cmdKeysList)),
	)
	return h.Run(ctx, options...)
}

// DeleteKey deletes a key from Hermes keys.
func (h *Hermes) DeleteKey(ctx context.Context, chainID, keyName string, options ...Option) error {
	options = append(
		options,
		WithFlags(Flags{
			FlagChain:   chainID,
			FlagKeyName: keyName,
		}),
		WithArgs(string(cmdKeys), string(cmdKeysDelete)),
	)
	return h.Run(ctx, options...)
}

// CreateClient creates a new relayer client.
func (h *Hermes) CreateClient(ctx context.Context, hostChain, referenceChain string, options ...Option) error {
	options = append(options, WithFlags(
		Flags{
			FlagHostChain:      hostChain,
			FlagReferenceChain: referenceChain,
		},
	))
	options = append(options, WithArgs(string(cmdCreate), string(cmdClient)))
	return h.Run(ctx, options...)
}

// CreateConnection creates a new relayer connection.
func (h *Hermes) CreateConnection(ctx context.Context, chainA, clientA, clientB string, options ...Option) error {
	options = append(options, WithFlags(
		Flags{
			FlagChainA:  chainA,
			FlagClientA: clientA,
			FlagClientB: clientB,
		},
	))
	options = append(options, WithArgs(string(cmdCreate), string(cmdConnection)))
	return h.Run(ctx, options...)
}

// CreateChannel creates a new relayer channel.
func (h *Hermes) CreateChannel(ctx context.Context, chainA, connA, portA, portB string, options ...Option) error {
	options = append(options, WithFlags(
		Flags{
			FlagChainA:      chainA,
			FlagConnectionA: connA,
			FlagPortA:       portA,
			FlagPortB:       portB,
		},
	))
	options = append(options, WithArgs(string(cmdCreate), string(cmdChannel)))
	return h.Run(ctx, options...)
}

// QueryChannels query all Hermes channels based in a chain id.
func (h *Hermes) QueryChannels(ctx context.Context, showCounterparty bool, chain string, options ...Option) error {
	flags := Flags{
		FlagChain: chain,
	}
	if showCounterparty {
		flags[FlagShowCounterparty] = true
	}
	options = append(
		options,
		WithFlags(flags),
		WithArgs(string(cmdQuery), string(cmdChannels)),
	)
	return h.Run(ctx, options...)
}

// Start starts the Hermes relayer.
func (h *Hermes) Start(ctx context.Context, options ...Option) error {
	options = append(options, WithArgs(string(cmdStart)))
	return h.Run(ctx, options...)
}

// Run runs a Hermes command using the options.
func (h *Hermes) Run(ctx context.Context, options ...Option) error {
	c := configs{}
	for _, o := range options {
		o(&c)
	}

	cmd := []string{h.path}

	// the config and json flag should be added before the hermes subcommands
	if c.config != "" {
		cmd = append(cmd, fmt.Sprintf("--%s=%s", FlagConfig, c.config))
	}

	cmd = append(cmd, c.args...)

	for flag, value := range c.flags {
		if v, ok := value.(bool); ok && v {
			cmd = append(cmd, fmt.Sprintf("--%s", flag))
		} else {
			cmd = append(cmd, fmt.Sprintf("--%s=%s", flag, value))
		}
	}

	stdout := c.stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	stderr := c.stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	// Hermes returns JSON formatted errors to stdout when something fails during execution.
	// A secondary buffer is used to be able to read the error output and to allow the caller
	// to also read the output from the configured stdout writer later on.
	var out bytes.Buffer
	stdout = io.MultiWriter(stdout, &out)
	err := exec.Exec(ctx, cmd, exec.StepOption(step.Stdout(stdout)), exec.StepOption(step.Stderr(stderr)))
	if err != nil {
		// Try to parse stdout as a Hermes formatted error
		if err := parseErrFromOutput(out.Bytes()); err != nil {
			return err
		}
		// Otherwise return the execution error
		return err
	}
	return nil
}

// UnmarshalResult unmarshal the command result into a interface.
func UnmarshalResult(data []byte, v any) error {
	var r Result
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	if r.Status != ResultSuccess {
		return fmt.Errorf("error result (%T) error: %v", v, string(r.Result))
	}
	return json.Unmarshal(r.Result, v)
}

// ValidateResult validate if the cmd result is success.
func ValidateResult(data []byte) error {
	var r Result
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	if r.Status != ResultSuccess {
		return fmt.Errorf("%w: %v", ErrResult, string(r.Result))
	}
	return nil
}

// parseErrFromOutput parses any error found in the Hermes output.
// Error are sent to stdout as JSON lines where the last line might
// contain the final error message returned by Hermes. Previous lines
// might contain standard logging entries.
func parseErrFromOutput(out []byte) error {
	out = bytes.TrimSpace(out)
	if len(out) > 0 {
		lines := bytes.Split(out, []byte("\n"))
		err := ValidateResult(lines[len(lines)-1])
		if errors.Is(err, ErrResult) {
			return err
		}
	}
	return nil
}
