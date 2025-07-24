package ssh

import (
	"context"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
)

// OS returns the operating system type of the remote server.
func (s *SSH) OS(ctx context.Context) (string, error) {
	v, err := s.Uname(ctx)
	if err != nil {
		return "", err
	}
	return strings.ToLower(v), nil
}

// Arch returns the architecture type of the remote server.
func (s *SSH) Arch(ctx context.Context) (string, error) {
	v, err := s.Uname(ctx, "-m")
	if err != nil {
		return "", err
	}
	if arch, ok := archMap[v]; ok {
		v = arch
	}
	return strings.ToLower(v), nil
}

// Target returns the build target for the remote server based on its OS and architecture.
func (s *SSH) Target(ctx context.Context) (string, error) {
	osType, err := s.OS(ctx)
	if err != nil {
		return "", err
	}

	arch, err := s.Arch(ctx)
	if err != nil {
		return "", err
	}

	return gocmd.BuildTarget(osType, arch), nil
}

// Uname runs the "uname" command with the specified arguments on the remote server.
func (s *SSH) Uname(ctx context.Context, args ...string) (string, error) {
	return s.RunCommand(ctx, "uname", args...)
}
