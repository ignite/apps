// Package ssh provides functionalities for establishing SSH connections
// and performing various operations such as file uploads, command execution,
// and managing remote environments.

package ssh

import (
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/gocmd"
	"github.com/ignite/cli/v28/ignite/pkg/randstr"
	"github.com/melbahja/goph"
	"github.com/pkg/sftp"
)

const (
	workdir = "spaceship"
)

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

// Workspace returns the workspace directory for the SSH session.
func (s *SSH) Workspace() string {
	return filepath.Join(workdir, s.workspace)
}

// Bin returns the binary directory within the workspace.
func (s *SSH) Bin() string {
	return filepath.Join(s.Workspace(), "bin")
}

// Home returns the home directory within the workspace.
func (s *SSH) Home() string {
	return filepath.Join(s.Workspace(), "home")
}

// Genesis returns the path to the genesis.json file within the home directory.
func (s *SSH) Genesis() string {
	return filepath.Join(s.Home(), "config", "genesis.json")
}

// Log returns the log directory within the workspace.
func (s *SSH) Log() string {
	return filepath.Join(s.Workspace(), "log")
}

// RunnerScript returns the path to the runner script within the workspace.
func (s *SSH) RunnerScript() string {
	return filepath.Join(s.Workspace(), "run.sh")
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
	if err != nil {
		return errors.Wrapf(err, "Failed to connect to %v", s)
	}

	s.sftpClient, err = s.client.NewSftp()
	if err != nil {
		return err
	}

	return s.ensureEnvironment()
}

// ensureEnvironment ensures that the necessary directories exist on the remote server.
func (s *SSH) ensureEnvironment() error {
	if err := s.sftpClient.MkdirAll(s.Bin()); err != nil {
		return errors.Wrapf(err, "failed to create bin dir %s", s.Bin())
	}
	if err := s.sftpClient.MkdirAll(s.Home()); err != nil {
		return errors.Wrapf(err, "failed to create home dir %s", s.Home())
	}
	if err := s.sftpClient.MkdirAll(s.Log()); err != nil {
		return errors.Wrapf(err, "failed to create home dir %s", s.Home())
	}
	return nil
}

// ensureLocalBin uploads the specified binary to the remote server's bin directory.
func (s *SSH) ensureLocalBin(name string, progressCallback ProgressCallback) error {
	// find ignite binary path
	path, err := exec.LookPath(name)
	if err != nil {
		return err
	}
	_, err = s.UploadBinary(path, progressCallback)
	if err != nil {
		return err
	}
	return nil
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
	return s.RunCommand(ctx, s.RunnerScript(), args...)
}

// HasGenesis checks if the genesis file exists on the remote server.
func (s *SSH) HasGenesis(ctx context.Context) bool {
	return s.FileExist(ctx, s.Genesis())
}

// HasRunnerScript checks if the runner script file exists on the remote server.
func (s *SSH) HasRunnerScript(ctx context.Context) bool {
	return s.FileExist(ctx, s.RunnerScript())
}

// FolderExist checks if a directory exists at the specified path on the remote server.
// It returns true if the directory exists, otherwise false.
func (s *SSH) FolderExist(ctx context.Context, path string) bool {
	return s.exist(ctx, path, false)
}

// FileExist checks if a file exists at the specified path on the remote server.
// It returns true if the file exists, otherwise false.
func (s *SSH) FileExist(ctx context.Context, path string) bool {
	return s.exist(ctx, path, true)
}

// exist checks if a file or directory exists at the specified path on the remote server.
// If isFile is true, it checks for a file, otherwise it checks for a directory.
// It returns true if the specified file or directory exists, otherwise false.
func (s *SSH) exist(ctx context.Context, path string, isFile bool) bool {
	cmd := fmt.Sprintf("[ -d '%s' ] && echo 'true'", path)
	if isFile {
		cmd = fmt.Sprintf("[ -f '%s' ] && echo 'true'", path)
	}
	exist, err := s.RunCommand(ctx, cmd)
	if err != nil {
		return false
	}
	return exist == "true"
}

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
	os, err := s.OS(ctx)
	if err != nil {
		return "", err
	}

	arch, err := s.Arch(ctx)
	if err != nil {
		return "", err
	}

	return gocmd.BuildTarget(os, arch), nil
}

// Uname runs the "uname" command with the specified arguments on the remote server.
func (s *SSH) Uname(ctx context.Context, args ...string) (string, error) {
	return s.RunCommand(ctx, "uname", args...)
}
