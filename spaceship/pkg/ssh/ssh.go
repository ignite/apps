// Package ssh provides functionalities for establishing SSH connections
// and performing various operations such as file uploads, command execution,
// and managing remote environments.
package ssh

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/randstr"
	"github.com/manifoldco/promptui"
	"github.com/melbahja/goph"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh/knownhosts"
)

const workdir = "spaceship"

// SSH represents the SSH configuration and clients for connecting and interacting
// with remote servers via SSH.
type SSH struct {
	username    string
	password    string
	host        string
	port        string
	rawKey      string
	key         string
	keyPassword string
	workspace   string
	client      *goph.Client
	sftpClient  *sftp.Client
}

// Option configures SSH settings.
type Option func(*SSH) error

// WithUser sets the SSH username.
func WithUser(username string) Option {
	return func(o *SSH) error {
		if o.username != "" {
			return nil
		}
		o.username = strings.TrimSpace(username)
		return nil
	}
}

// WithPassword sets the SSH password.
func WithPassword(password string) Option {
	return func(o *SSH) error {
		if o.password != "" {
			return nil
		}
		o.password = strings.TrimSpace(password)
		return nil
	}
}

// WithPort sets the SSH port.
func WithPort(port string) Option {
	return func(o *SSH) error {
		if o.port != "" {
			return nil
		}
		o.port = strings.TrimSpace(port)
		return nil
	}
}

// WithRawKey sets the SSH raw key.
func WithRawKey(rawKey string) Option {
	return func(o *SSH) error {
		if o.rawKey != "" {
			return nil
		}
		o.rawKey = strings.TrimSpace(rawKey)
		return nil
	}
}

// WithKey sets the SSH key.
func WithKey(key string) Option {
	return func(o *SSH) error {
		if o.key != "" {
			return nil
		}
		o.key = strings.TrimSpace(key)
		return nil
	}
}

// WithKeyPassword sets the SSH key password.
func WithKeyPassword(keyPassword string) Option {
	return func(o *SSH) error {
		if o.keyPassword != "" {
			return nil
		}
		o.keyPassword = strings.TrimSpace(keyPassword)
		return nil
	}
}

// WithWorkspace sets the SSH workspace.
func WithWorkspace(workspace string) Option {
	return func(o *SSH) error {
		o.workspace = strings.TrimSpace(workspace)
		return nil
	}
}

// New creates a new SSH object with the given host and options.
func New(host string, options ...Option) (*SSH, error) {
	host, port, username, password, err := parseURI(host)
	if err != nil {
		return nil, err
	}
	s := &SSH{
		username:  username,
		host:      host,
		port:      port,
		password:  password,
		workspace: randstr.Runes(10),
	}
	for _, apply := range options {
		if err := apply(s); err != nil {
			return nil, err
		}
	}
	return s, s.validate()
}

// parseURI parses the SSH URI and extracts the host, port, username, and password.
func parseURI(uri string) (host string, port string, username string, password string, err error) {
	uri = strings.TrimSpace(uri)
	if !strings.HasPrefix(uri, "ssh://") {
		uri = "ssh://" + uri
	}
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return "", "", "", "", errors.Wrapf(err, "error parsing URI %s", uri)
	}
	host = parsedURL.Hostname()
	port = parsedURL.Port()
	if port == "" {
		port = "22"
	}

	if parsedURL.User != nil {
		username = parsedURL.User.Username()
		password, _ = parsedURL.User.Password()
	}
	if username == "" {
		username = "root"
	}
	return host, port, username, password, nil
}

// validate checks if the SSH configuration is valid.
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

// auth returns the appropriate authentication method based on the SSH configuration.
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

// ensureEnvironment ensures that the necessary directories exist on the remote server.
func (s *SSH) ensureEnvironment() error {
	if err := s.sftpClient.MkdirAll(s.bin()); err != nil {
		return errors.Wrapf(err, "failed to create bin dir %s", s.bin())
	}
	if err := s.sftpClient.MkdirAll(s.Home()); err != nil {
		return errors.Wrapf(err, "failed to create home dir %s", s.Home())
	}
	if err := s.sftpClient.MkdirAll(s.Log()); err != nil {
		return errors.Wrapf(err, "failed to create home dir %s", s.Home())
	}
	return nil
}

// bin returns the binary directory within the workspace.
func (s *SSH) bin() string {
	return filepath.Join(s.Workspace(), "bin")
}

// Workspace returns the workspace directory for the SSH session.
func (s *SSH) Workspace() string {
	return filepath.Join(workdir, s.workspace)
}

// Close closes the SSH and SFTP clients.
func (s *SSH) Close() error {
	if err := s.sftpClient.Close(); err != nil {
		return err
	}
	return s.client.Close()
}

// Connect establishes the SSH connection and initializes the SFTP client.
func (s *SSH) Connect() error {
	auth, err := s.auth()
	if err != nil {
		return err
	}

	s.client, err = goph.New(s.username, s.host, auth)
	if errors.Is(err, &knownhosts.KeyError{}) {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Unknown host: %s. Do you want to proceed with the connection anyway", s.host),
			IsConfirm: true,
			Stdout:    os.Stdout,
			Stdin:     os.Stdin,
		}
		if _, err := prompt.Run(); err != nil {
			return err
		}
		s.client, err = goph.NewUnknown(s.username, s.host, auth)
	}
	if err != nil {
		return errors.Wrapf(err, "Failed to connect to %v", s)
	}

	s.sftpClient, err = s.client.NewSftp()
	if err != nil {
		return err
	}

	return s.ensureEnvironment()
}

// RunCommand runs a command on the remote server and returns the output.
func (s *SSH) RunCommand(ctx context.Context, name string, args ...string) (string, error) {
	cmd, err := s.client.CommandContext(ctx, name, args...)
	if err != nil {
		return "", err
	}
	cmdOut, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(cmdOut))
	if err != nil {
		return "", errors.Errorf(
			"%s: failed to run %s %s\n%s",
			err.Error(),
			name,
			strings.Join(args, " "),
			output,
		)
	}
	return output, nil
}
