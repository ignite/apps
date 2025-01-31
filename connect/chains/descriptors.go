package chains

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"path"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"

	"github.com/ignite/cli/v28/ignite/pkg/chainregistry"
)

type Conn struct {
	chain     chainregistry.Chain
	config    *ChainConfig
	configDir string

	client        *grpc.ClientConn
	protoFiles    *protoregistry.Files
	moduleOptions map[string]*autocliv1.ModuleOptions
}

func NewConn(chain chainregistry.Chain, cfg *ChainConfig) (*Conn, error) {
	configDir, err := configDir()
	if err != nil {
		return nil, err
	}

	return &Conn{
		chain:     chain,
		config:    cfg,
		configDir: configDir,
	}, nil
}

// fdsCacheFilename returns the filename for the cached file descriptor set.
func (c *Conn) fdsCacheFilename() string {
	return path.Join(c.configDir, fmt.Sprintf("%s.fds", c.chain.ChainName))
}

// appOptsCacheFilename returns the filename for the app options cache file.
func (c *Conn) appOptsCacheFilename() string {
	return path.Join(c.configDir, fmt.Sprintf("%s.autocli", c.chain.ChainName))
}

func (c *Conn) Load(ctx context.Context) error {
	var err error
	fdSet := &descriptorpb.FileDescriptorSet{}
	fdsFilename := c.fdsCacheFilename()

	if _, err := os.Stat(fdsFilename); os.IsNotExist(err) {
		client, err := c.Connect()
		if err != nil {
			return err
		}

		reflectionClient := reflectionv1.NewReflectionServiceClient(client)
		fdRes, err := reflectionClient.FileDescriptors(ctx, &reflectionv1.FileDescriptorsRequest{})
		if err != nil {
			return fmt.Errorf("error getting file descriptors: %w, this chain is using a too old version of the Cosmos SDK", err)
		}

		fdSet = &descriptorpb.FileDescriptorSet{File: fdRes.Files}
		bz, err := proto.Marshal(fdSet)
		if err != nil {
			return err
		}

		if err = os.WriteFile(fdsFilename, bz, 0o600); err != nil {
			return err
		}
	} else {
		bz, err := os.ReadFile(fdsFilename)
		if err != nil {
			return err
		}

		if err = proto.Unmarshal(bz, fdSet); err != nil {
			return err
		}
	}

	c.protoFiles, err = protodesc.FileOptions{AllowUnresolvable: true}.NewFiles(fdSet)
	if err != nil {
		return fmt.Errorf("error building protoregistry.Files: %w", err)
	}

	appOptsFilename := c.appOptsCacheFilename()
	if _, err := os.Stat(appOptsFilename); os.IsNotExist(err) {
		client, err := c.Connect()
		if err != nil {
			return err
		}

		autocliQueryClient := autocliv1.NewQueryClient(client)
		appOptsRes, err := autocliQueryClient.AppOptions(ctx, &autocliv1.AppOptionsRequest{})
		if err != nil {
			return fmt.Errorf("error getting autocli config: %w, this chain is using a too old version of the Cosmos SDK", err)
		}

		bz, err := proto.Marshal(appOptsRes)
		if err != nil {
			return err
		}

		if err := os.WriteFile(appOptsFilename, bz, 0o600); err != nil {
			return err
		}

		c.moduleOptions = appOptsRes.ModuleOptions
	} else {
		bz, err := os.ReadFile(appOptsFilename)
		if err != nil {
			return err
		}

		var appOptsRes autocliv1.AppOptionsResponse
		if err := proto.Unmarshal(bz, &appOptsRes); err != nil {
			return err
		}

		c.moduleOptions = appOptsRes.ModuleOptions
	}

	return nil
}

func (c *Conn) Connect() (*grpc.ClientConn, error) {
	if c.client != nil {
		return c.client, nil
	}

	var err error
	creds := credentials.NewTLS(&tls.Config{
		MinVersion: tls.VersionTLS12,
	})

	// we use Dial here instead of NewClient as we want to attempt to connect to the gRPC server immediately
	// and fallback to another connection if it fails
	c.client, err = grpc.Dial(c.config.GRPCEndpoint, grpc.WithTransportCredentials(creds)) //nolint:staticcheck: we want to use dial
	if err != nil {
		creds = insecure.NewCredentials()
		c.client, err = grpc.NewClient(c.config.GRPCEndpoint, grpc.WithTransportCredentials(creds))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return c.client, nil
}
