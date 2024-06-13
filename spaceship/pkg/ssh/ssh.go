package ssh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

const chainAppName = "chain-app"

type SSH struct {
	username   string
	password   string
	host       string
	port       string
	binaryPath string
	client     *ssh.Client
}

// NewFromURI creates a new ssh object from a URI URL.
func NewFromURI(uri, binaryPath string) (*SSH, error) {
	// parse the URI.
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "Error parsing URL %s", uri)
	}
	var (
		host     = parsedURL.Hostname()
		port     = parsedURL.Port()
		username = ""
		password = ""
	)

	// extract user information.
	if parsedURL.User != nil {
		username = parsedURL.User.Username()
		password, _ = parsedURL.User.Password()
	}
	return New(username, password, host, port, binaryPath), nil
}

// New creates a new ssh object.
func New(username, password, host, port, binaryPath string) *SSH {
	return &SSH{
		username:   username,
		password:   password,
		host:       host,
		port:       port,
		binaryPath: binaryPath,
	}
}

// Close closes the SSH client.
func (s *SSH) Close() error {
	return s.client.Close()
}

// Connect connects the SSH client.
func (s *SSH) Connect() (err error) {
	// SSH connection setup
	config := &ssh.ClientConfig{
		User:            s.username,
		Auth:            []ssh.AuthMethod{ssh.Password(s.password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// Establish SSH connection
	s.client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", s.host, s.port), config)
	if err != nil {
		return errors.Wrapf(err, "Failed to connect to %s:%s", s.host, s.port)
	}
	return nil
}

func (s *SSH) NewSession() (*ssh.Session, error) {
	return s.client.NewSession()
}

func (s *SSH) RunApp() error {
	// create a new session for running the application.
	runSession, err := s.NewSession()
	if err != nil {
		return errors.Wrap(err, "failed to create SSH session for running app")
	}
	defer runSession.Close()

	// run the application on the remote server.
	if err := runSession.Run(fmt.Sprintf("./%s", chainAppName)); err != nil {
		return errors.Wrap(err, "failed to run app")
	}
	return nil
}

func (s *SSH) SendBinary(ctx context.Context) error {
	// create a new SSH session.
	session, err := s.NewSession()
	if err != nil {
		return errors.Wrap(err, "failed to create SSH session")
	}
	defer session.Close()

	// read the built application.
	appData, err := os.ReadFile(filepath.Join(filepath.Dir(s.binaryPath), chainAppName))
	if err != nil {
		log.Fatalf("Failed to read built app: %s", err)
	}

	// start SCP process.
	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error {
		w, err := session.StdinPipe()
		if err != nil {
			return err
		}
		defer w.Close()
		if _, err := fmt.Fprintln(w, "C0755", len(appData), chainAppName); err != nil {
			return err
		}
		if _, err := w.Write(appData); err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, "\x00")
		return err
	})
	if err := errg.Wait(); err != nil {
		return err
	}

	if err := session.Run("/usr/bin/scp -tr ."); err != nil {
		return errors.Wrap(err, "Failed to transfer applo")
	}
	return nil
}
