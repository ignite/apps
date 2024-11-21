package faucet

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/blang/semver/v4"
	"github.com/ignite/cli/v28/ignite/config"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/xfilepath"
	"github.com/ignite/cli/v28/ignite/pkg/xos"

	"github.com/ignite/apps/spaceship/pkg/tarball"
)

const (
	faucetBinaryName     = "faucet"
	binaryCacheDirectory = "spaceship/bin"
	faucetLastRelease    = "https://github.com/ignite/faucet/releases/latest/download"
)

var faucetVersion = semver.MustParse("0.0.1")

func faucetReleaseName(target string) string {
	return fmt.Sprintf("%s/faucet_%s_%s.tar.gz", faucetLastRelease, faucetVersion.String(), target)
}

func FetchBinary(ctx context.Context, target string) (string, error) {
	binPath, err := binCachePath()
	if err != nil {
		return "", err
	}

	// Check if the binary already exists in the ignite cache.
	if _, err := os.Stat(binPath); err == nil {
		return binPath, nil
	}

	// Create a temporary folder to extract the faucet binary.
	tempDir, err := os.MkdirTemp("", "faucet")
	if err != nil {
		return "", errors.Errorf("failed to create temp dir: %w", err)
	}

	binaryURL := faucetReleaseName(target)

	// Download the binary.
	resp, err := http.Get(binaryURL)
	if err != nil {
		return "", errors.Errorf("failed to download faucet binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("failed to fetch faucet binary: %s status", resp.Status)
	}

	// Extract the binary tarball.
	extracted, err := tarball.ExtractData(ctx, resp.Body, tempDir, faucetBinaryName)
	if err != nil {
		return "", err
	}
	if len(extracted) == 0 {
		return "", errors.Errorf("zero files extracted from %s faucet the tarball: %s", target, binaryURL)
	}

	return binPath, xos.Rename(extracted[0], binPath)
}

func BinaryName() string {
	return fmt.Sprintf("%s_%s", faucetBinaryName, faucetVersion.String())
}

func binCacheDirPath() (string, error) {
	return xfilepath.Join(config.DirPath, xfilepath.Path(binaryCacheDirectory))()
}

func binCachePath() (string, error) {
	dirPath, err := binCacheDirPath()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dirPath, BinaryName())
	return path, os.MkdirAll(dirPath, 0o755)
}
