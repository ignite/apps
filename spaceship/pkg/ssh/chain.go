package ssh

import (
	"context"
	"path/filepath"
)

// genesis returns the path to the genesis.json file within the home directory.
func (s *SSH) genesis() string {
	return filepath.Join(s.Home(), "config", "genesis.json")
}

// runnerScript returns the path to the runner script within the workspace.
func (s *SSH) runnerScript() string {
	return filepath.Join(s.Workspace(), "run.sh")
}

// Home returns the home directory within the workspace.
func (s *SSH) Home() string {
	return filepath.Join(s.Workspace(), "home")
}

// HasGenesis checks if the genesis file exists on the remote server.
func (s *SSH) HasGenesis(ctx context.Context) bool {
	return s.FileExist(ctx, s.genesis())
}

// HasRunnerScript checks if the runner script file exists on the remote server.
func (s *SSH) HasRunnerScript(ctx context.Context) bool {
	return s.FileExist(ctx, s.runnerScript())
}

// Start runs the "start" script on the remote server.
func (s *SSH) Start(ctx context.Context) (string, error) {
	return s.runScript(ctx, "start")
}

// Restart runs the "restart" script on the remote server.
func (s *SSH) Restart(ctx context.Context) (string, error) {
	return s.runScript(ctx, "restart")
}

// Stop runs the "stop" script on the remote server.
func (s *SSH) Stop(ctx context.Context) (string, error) {
	return s.runScript(ctx, "stop")
}

// Status runs the "status" script on the remote server.
func (s *SSH) Status(ctx context.Context) (string, error) {
	return s.runScript(ctx, "status")
}

// runScript runs the specified script with arguments on the remote server.
func (s *SSH) runScript(ctx context.Context, args ...string) (string, error) {
	return s.RunCommand(ctx, s.runnerScript(), args...)
}
