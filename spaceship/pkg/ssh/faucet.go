package ssh

import (
	"context"
	"fmt"
	"github.com/blang/semver/v4"
	"github.com/ignite/apps/spaceship/pkg/tarball"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"net/http"
	"os"
)

const (
	faucetBinName     = "faucet"
	faucetLastRelease = "https://github.com/ignite/faucet/releases/latest/download"
)

var faucetVersion = semver.MustParse("0.0.1")

func faucetReleaseName(target string) string {
	return fmt.Sprintf("%s/faucet_%s_%s.tar.gz", faucetLastRelease, faucetVersion.String(), target)
}

func fetchFaucetBinary(ctx context.Context, target string) (string, error) {
	tempDir, err := os.MkdirTemp("", "faucet")
	if err != nil {
		return "", errors.Errorf("failed to create temp dir: %w", err)
	}

	binaryURL := faucetReleaseName(target)
	resp, err := http.Get(binaryURL)
	if err != nil {
		return "", errors.Errorf("failed to download faucet binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("failed to fetch faucet binary: %s status", resp.Status)
	}

	extracted, err := tarball.ExtractData(ctx, resp.Body, tempDir, faucetBinName)
	if err != nil {
		return "", err
	}
	if len(extracted) == 0 {
		return "", errors.Errorf("zero files extracted from %s faucet the tarball: %s", target, binaryURL)
	}

	return "", nil
}
