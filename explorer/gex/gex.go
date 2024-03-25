package gex

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"os"

	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/localfs"
	"github.com/ignite/ignite-files/gex"
)

type (
	// Gex represents the Gex binary structure.
	Gex struct {
		path    string
		host    string
		port    string
		ssl     bool
		stdout  io.Writer
		stderr  io.Writer
		stdin   io.Reader
		cleanup func()
	}
	// Option configures the gex options.
	Option func(*Gex)
)

// newGex returns a Gex with default options.
func newGex() *Gex {
	return &Gex{
		host:    "localhost",
		port:    "26657",
		ssl:     false,
		stdout:  os.Stdout,
		stderr:  os.Stderr,
		stdin:   os.Stdin,
		cleanup: nil,
	}
}

// WithHost set the gex host.
func WithHost(host string) Option {
	return func(m *Gex) {
		m.host = host
	}
}

// WithPort set the gex port.
func WithPort(port string) Option {
	return func(m *Gex) {
		m.port = port
	}
}

// WithSSL set gex SSL.
func WithSSL(ssl bool) Option {
	return func(m *Gex) {
		m.ssl = ssl
	}
}

// WithStdout set gex Stdout.
func WithStdout(stdout io.Writer) Option {
	return func(m *Gex) {
		m.stdout = stdout
	}
}

// WithStdErr set gex StdErr.
func WithStdErr(stderr io.Writer) Option {
	return func(m *Gex) {
		m.stderr = stderr
	}
}

// WithStdIn set gex StdIn.
func WithStdIn(stdin io.Reader) Option {
	return func(m *Gex) {
		m.stdin = stdin
	}
}

// New returns the Gex binary executable.
func New(options ...Option) (*Gex, error) {
	g := newGex()
	for _, apply := range options {
		apply(g)
	}

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

	g.path, g.cleanup, err = localfs.SaveBytesTemp(binary, "gex", 0o755)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save gex binary as temp file")
	}

	return g, nil
}

// Cleanup clean the temporary Gex binary.
func (g *Gex) Cleanup() error {
	g.cleanup()
	return os.RemoveAll(g.path)
}

// Run runs gex with provided parameters.
func (g *Gex) Run(ctx context.Context) error {
	cmd := []string{
		g.path,
		"-h", g.host,
		"-p", g.port,
	}
	if g.ssl {
		cmd = append(cmd, "-s")
	}
	return exec.Exec(
		ctx,
		cmd,
		exec.StepOption(step.Stdout(g.stdout)),
		exec.StepOption(step.Stderr(g.stderr)),
		exec.StepOption(step.Stdin(g.stdin)),
	)
}
