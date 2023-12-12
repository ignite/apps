package gex

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"os"

	"github.com/ignite/ignite-files/gex"
	"github.com/pkg/errors"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/localfs"
)

// Gex represents the gex binary structure.
type Gex struct {
	path    string
	cleanup func()
}

// New returns the hermes binary executable.
func New() (*Gex, error) {
	// untar the binary.
	gzr, err := gzip.NewReader(bytes.NewReader(gex.Binary()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create gzip reader from stored binary")
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	if _, err := tr.Next(); err != nil {
		return nil, errors.Wrap(err, "failed to fetch next tar entry")
	}

	binary, err := io.ReadAll(tr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read tar entry")
	}

	path, cleanup, err := localfs.SaveBytesTemp(binary, "gex", 0o755)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save gex binary as temp file")
	}

	return &Gex{
		path:    path,
		cleanup: cleanup,
	}, nil
}

// Cleanup clean the temporary Gex binary.
func (g *Gex) Cleanup() error {
	g.cleanup()
	return os.RemoveAll(g.path)
}

// Run runs gex with provided parameters.
func (g *Gex) Run(ctx context.Context, stdout, stderr io.Writer, host, port string, ssl bool) error {
	cmd := []string{g.path}

	if host != "" {
		cmd = append(cmd, "-h", host)
	}
	if port != "" {
		cmd = append(cmd, "-p", port)
	}
	if ssl {
		cmd = append(cmd, "-s")
	}

	return exec.Exec(ctx, cmd, exec.StepOption(step.Stdout(stdout)), exec.StepOption(step.Stderr(stderr)))
}
