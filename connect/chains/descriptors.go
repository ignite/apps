package chains

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"path"

	authv1betav1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Conn struct {
	chainName string
	config    *ChainConfig
	configDir string
	client    *grpc.ClientConn

	ProtoFiles    *protoregistry.Files
	ModuleOptions map[string]*autocliv1.ModuleOptions
}

func NewConn(chainName string, cfg *ChainConfig) (*Conn, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return nil, err
	}

	return &Conn{
		chainName: chainName,
		config:    cfg,
		configDir: configDir,
	}, nil
}

// fdsCacheFilename returns the filename for the cached file descriptor set.
func (c *Conn) fdsCacheFilename() string {
	return path.Join(c.configDir, fmt.Sprintf("%s.fds", c.chainName))
}

// appOptsCacheFilename returns the filename for the app options cache file.
func (c *Conn) appOptsCacheFilename() string {
	return path.Join(c.configDir, fmt.Sprintf("%s.autocli", c.chainName))
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
			return fmt.Errorf("error getting file descriptors: %w", err)
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

	c.ProtoFiles, err = protodesc.FileOptions{AllowUnresolvable: true}.NewFiles(fdSet)
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
			return fmt.Errorf("error getting autocli config: %w", err)
		}

		bz, err := proto.Marshal(appOptsRes)
		if err != nil {
			return err
		}

		if err := os.WriteFile(appOptsFilename, bz, 0o600); err != nil {
			return err
		}

		c.ModuleOptions = appOptsRes.ModuleOptions
	} else {
		bz, err := os.ReadFile(appOptsFilename)
		if err != nil {
			return err
		}

		var appOptsRes autocliv1.AppOptionsResponse
		if err := proto.Unmarshal(bz, &appOptsRes); err != nil {
			return err
		}

		c.ModuleOptions = appOptsRes.ModuleOptions
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

	c.client, err = grpc.NewClient(c.config.GRPCEndpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	// try connection by querying an endpoint
	// fallback to insecure if it doesn't work
	authClient := authv1betav1.NewQueryClient(c.client)
	if _, err = authClient.Params(context.Background(), &authv1betav1.QueryParamsRequest{}); err != nil {
		creds = insecure.NewCredentials()
		c.client, err = grpc.NewClient(c.config.GRPCEndpoint, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
		}
	}

	return c.client, nil
}
