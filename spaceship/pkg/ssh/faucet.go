package ssh

import (
	"context"
	"path/filepath"
	"strconv"

	"github.com/ignite/apps/spaceship/pkg/faucet"
)

type faucetOpts struct {
	port            uint64
	keyringBackend  string
	sdkVersion      string
	accountName     string
	mnemonic        string
	keyringPassword string
	cliName         string
	denoms          string
	creditAmount    string
	maxCredit       string
	feeAmount       string
	node            string
	coinType        string
	home            string
}

// FaucetOption configures faucet options.
type FaucetOption func(*faucetOpts)

// WithFaucetPort sets the faucet port.
func WithFaucetPort(port uint64) FaucetOption {
	return func(o *faucetOpts) {
		o.port = port
	}
}

// WithFaucetSdkVersion sets the faucet sdk version.
func WithFaucetSdkVersion(sdkVersion string) FaucetOption {
	return func(o *faucetOpts) {
		o.sdkVersion = sdkVersion
	}
}

// WithFaucetAccountName sets the faucet account name.
func WithFaucetAccountName(accountName string) FaucetOption {
	return func(o *faucetOpts) {
		o.accountName = accountName
	}
}

// WithFaucetMnemonic sets the faucet mnemonic.
func WithFaucetMnemonic(mnemonic string) FaucetOption {
	return func(o *faucetOpts) {
		o.mnemonic = mnemonic
	}
}

// WithFaucetKeyringPassword sets the faucet keyring password.
func WithFaucetKeyringPassword(keyringPassword string) FaucetOption {
	return func(o *faucetOpts) {
		o.keyringPassword = keyringPassword
	}
}

// WithFaucetCliName sets the faucet CLI name.
func WithFaucetCliName(cliName string) FaucetOption {
	return func(o *faucetOpts) {
		o.cliName = cliName
	}
}

// WithFaucetDenoms sets the faucet denoms.
func WithFaucetDenoms(denoms string) FaucetOption {
	return func(o *faucetOpts) {
		o.denoms = denoms
	}
}

// WithFaucetCreditAmount sets the faucet credit amount.
func WithFaucetCreditAmount(creditAmount string) FaucetOption {
	return func(o *faucetOpts) {
		o.creditAmount = creditAmount
	}
}

// WithFaucetMaxCredit sets the faucet max credit.
func WithFaucetMaxCredit(maxCredit string) FaucetOption {
	return func(o *faucetOpts) {
		o.maxCredit = maxCredit
	}
}

// WithFaucetFeeAmount sets the faucet fee amount.
func WithFaucetFeeAmount(feeAmount string) FaucetOption {
	return func(o *faucetOpts) {
		o.feeAmount = feeAmount
	}
}

// WithFaucetNode sets the faucet node.
func WithFaucetNode(node string) FaucetOption {
	return func(o *faucetOpts) {
		o.node = node
	}
}

// WithFaucetCoinType sets the faucet coin type.
func WithFaucetCoinType(coinType string) FaucetOption {
	return func(o *faucetOpts) {
		o.coinType = coinType
	}
}

// WithFaucetHome sets the faucet home.
func WithFaucetHome(home string) FaucetOption {
	return func(o *faucetOpts) {
		o.home = home
	}
}

// faucet returns the path to the faucet script within the workspace.
func (s *SSH) faucet() string {
	return filepath.Join(s.bin(), faucet.BinaryName())
}

// RunFaucet runs the faucet on the remote server.
func (s *SSH) RunFaucet(
	ctx context.Context,
	options ...FaucetOption,
) (string, error) {
	args := make([]string, 0)
	o := &faucetOpts{}
	for _, apply := range options {
		apply(o)
	}
	if o.port != 0 {
		args = append(args, "--port", strconv.FormatUint(o.port, 10))
	}
	if o.keyringBackend != "" {
		args = append(args, "--keyring-backend", o.keyringBackend)
	}
	if o.sdkVersion != "" {
		args = append(args, "--sdk-version", o.sdkVersion)
	}
	if o.accountName != "" {
		args = append(args, "--account-name", o.accountName)
	}
	if o.mnemonic != "" {
		args = append(args, "--mnemonic", o.mnemonic)
	}
	if o.keyringPassword != "" {
		args = append(args, "--keyring-password", o.keyringPassword)
	}
	if o.cliName != "" {
		args = append(args, "--cli-name", o.cliName)
	}
	if o.denoms != "" {
		args = append(args, "--denoms", o.denoms)
	}
	if o.creditAmount != "" {
		args = append(args, "--credit-amount", o.creditAmount)
	}
	if o.maxCredit != "" {
		args = append(args, "--max-credit", o.maxCredit)
	}
	if o.feeAmount != "" {
		args = append(args, "--fee-amount", o.feeAmount)
	}
	if o.node != "" {
		args = append(args, "--node", o.node)
	}
	if o.coinType != "" {
		args = append(args, "--coin-type", o.coinType)
	}
	if o.home != "" {
		args = append(args, "--home", o.home)
	}
	return s.RunCommand(ctx, s.faucet(), args...)
}
