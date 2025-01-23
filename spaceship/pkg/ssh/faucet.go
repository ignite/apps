package ssh

import (
	"context"
	"path/filepath"
	"strconv"

	"github.com/ignite/apps/spaceship/pkg/faucet"
)

// faucet returns the path to the faucet script within the workspace.
func (s *SSH) faucet() string {
	return filepath.Join(s.Bin(), faucet.BinaryName())
}

// faucetScript returns the path to the faucet runner script within the workspace.
func (s *SSH) faucetScript() string {
	return filepath.Join(s.Workspace(), "faucet.sh")
}

// HasFaucetScript checks if the runner faucet script file exists on the remote server.
func (s *SSH) HasFaucetScript(ctx context.Context) bool {
	return s.FileExist(ctx, s.faucetScript())
}

// FaucetStart runs the faucet "start" script on the remote server.
func (s *SSH) FaucetStart(ctx context.Context, port uint64) (string, error) {
	return s.runFaucetScript(ctx, "start", strconv.FormatUint(port, 10))
}

// FaucetRestart runs the faucet "restart" script on the remote server.
func (s *SSH) FaucetRestart(ctx context.Context, port uint64) (string, error) {
	return s.runFaucetScript(ctx, "restart", strconv.FormatUint(port, 10))
}

// FaucetStop runs the faucet "stop" script on the remote server.
func (s *SSH) FaucetStop(ctx context.Context) (string, error) {
	return s.runFaucetScript(ctx, "stop")
}

// FaucetStatus runs the faucet "status" script on the remote server.
func (s *SSH) FaucetStatus(ctx context.Context) (string, error) {
	return s.runFaucetScript(ctx, "status")
}

// runFaucetScript runs the specified faucet script with arguments on the remote server.
func (s *SSH) runFaucetScript(ctx context.Context, args ...string) (string, error) {
	return s.RunCommand(ctx, s.faucetScript(), args...)
}
