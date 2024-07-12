package ssh

import (
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/melbahja/goph"
	"github.com/pkg/sftp"
)

const (
	igniteAppName = "ignite"
	workdir       = "spaceship"
)

var (
	igniteWorkdir = filepath.Join(workdir, "ignite")
	igniteBinary  = filepath.Join(igniteWorkdir, "ignite")
)

type SSH struct {
	username    string
	password    string
	host        string
	port        string
	rawKey      string
	key         string
	keyPassword string
	client      *goph.Client
}

// Option configures ssh config.
type Option func(*SSH) error

// WithUser set SSH username.
func WithUser(username string) Option {
	return func(o *SSH) error {
		o.username = strings.TrimSpace(username)
		return nil
	}
}

// WithPassword set SSH password.
func WithPassword(password string) Option {
	return func(o *SSH) error {
		o.password = strings.TrimSpace(password)
		return nil
	}
}

// WithPort set SSH port.
func WithPort(port string) Option {
	return func(o *SSH) error {
		o.port = strings.TrimSpace(port)
		return nil
	}
}

// WithRawKey set SSH raw key.
func WithRawKey(rawKey string) Option {
	return func(o *SSH) error {
		o.rawKey = strings.TrimSpace(rawKey)
		return nil
	}
}

// WithKey set SSH key.
func WithKey(key string) Option {
	return func(o *SSH) error {
		o.key = strings.TrimSpace(key)
		return nil
	}
}

// WithKeyPassword set SSH key password.
func WithKeyPassword(keyPassword string) Option {
	return func(o *SSH) error {
		o.keyPassword = strings.TrimSpace(keyPassword)
		return nil
	}
}

// WithURI set SSH URI.
func WithURI(uri string) Option {
	return func(o *SSH) error {
		uri = strings.TrimSpace(uri)
		parsedURL, err := url.Parse(uri)
		if err != nil {
			return errors.Wrapf(err, "error parsing URI %s", uri)
		}
		o.host = parsedURL.Hostname()
		o.port = parsedURL.Port()
		o.username = ""
		o.password = ""

		// extract user information.
		if parsedURL.User != nil {
			o.username = parsedURL.User.Username()
			o.password, _ = parsedURL.User.Password()
		}
		return nil
	}
}

func (s *SSH) validate() error {
	switch {
	case s.username == "":
		return fmt.Errorf("ssh username is required")
	case s.key != "" && s.rawKey != "":
		return errors.New("ssh key and raw key are both set")
	case s.key != "" && s.password != "":
		return errors.New("ssh key and password are both set")
	case s.rawKey != "" && s.password != "":
		return errors.New("ssh raw key and password are both set")
	default:
		return nil
	}
}

func (s *SSH) auth() (goph.Auth, error) {
	switch {
	case s.rawKey != "":
		return goph.RawKey(s.rawKey, s.keyPassword)
	case s.key != "":
		return goph.Key(s.key, s.keyPassword)
	case s.password != "":
		return goph.Password(s.password), nil
	default:
		return goph.KeyboardInteractive(s.password), nil
	}
}

// New creates a new ssh object.
func New(host string, options ...Option) (*SSH, error) {
	s := &SSH{
		username: "root",
		host:     host,
		port:     "22",
	}
	for _, apply := range options {
		if err := apply(s); err != nil {
			return nil, err
		}
	}
	return s, s.validate()
}

// Close closes the SSH client.
func (s *SSH) Close() error {
	return s.client.Close()
}

// Connect connects the SSH client.
func (s *SSH) Connect(ctx context.Context) error {
	auth, err := s.auth()
	if err != nil {
		return err
	}
	s.client, err = goph.New(s.username, s.host, auth)
	if err != nil {
		return errors.Wrapf(err, "Failed to connect to %v", s)
	}

	return s.ensureEnvironment(ctx)
}

func (s *SSH) ensureEnvironment(ctx context.Context) error {
	sftp, err := s.client.NewSftp()
	if err != nil {
		return err
	}
	if err := s.ensureHomeFolder(sftp); err != nil {
		return err
	}
	if err := s.ensureIgniteBin(ctx, sftp); err != nil {
		return err
	}
	return nil
}

func (s *SSH) ensureHomeFolder(sftp *sftp.Client) error {
	if err := sftp.MkdirAll(igniteWorkdir); err != nil {
		return errors.Wrapf(err, "failed to create workdir %s", igniteWorkdir)
	}
	return nil
}

func (s *SSH) ensureIgniteBin(ctx context.Context, sftp *sftp.Client) error {
	if err := s.SendBinary(igniteAppName, igniteWorkdir); err != nil {
		return err
	}
	if err := sftp.Chmod(igniteBinary, 0o755); err != nil {
		return err
	}

	cmd, err := s.client.CommandContext(ctx, igniteBinary, "-h")
	if err != nil {
		return err
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	if len(out) == 0 {
		return errors.New("ignite binary doesn't exist")
	}
	return nil
}

func (s *SSH) SendFile(srcPath, dstPath string) error {
	srcPath, err := filepath.Abs(srcPath)
	if err != nil {
		return err
	}
	return s.client.Upload(srcPath, filepath.Join(dstPath, filepath.Base(srcPath)))
}

func (s *SSH) SendBinary(binaryName, dstPath string) error {
	path, err := exec.LookPath(binaryName)
	if err != nil {
		return err
	}
	if err := s.client.Upload(path, filepath.Join(dstPath, binaryName)); err != nil {
		return err
	}
	return nil
}
