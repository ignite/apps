package faucet

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/blang/semver/v4"
	"github.com/ignite/cli/v29/ignite/config"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xfilepath"
	"github.com/ignite/cli/v29/ignite/pkg/xos"

	"github.com/ignite/apps/spaceship/pkg/tarball"
)

const (
	faucetBinaryName     = "faucet"
	binaryCacheDirectory = "apps/spaceship/bin"
	faucetLastRelease    = "https://github.com/ignite/faucet/releases/download"
)

// faucetVersion specifies the current version of the faucet application.
var faucetVersion = semver.MustParse("0.0.3")

// faucetReleaseName constructs the download URL for a faucet binary tarball given the target platform.
func faucetReleaseName(target string) string {
	return fmt.Sprintf("%[1]v/v%[2]v/faucet_%[2]v_%[3]v.tar.gz", faucetLastRelease, faucetVersion.String(), target)
}

// FetchBinary downloads the faucet binary file from a specific target
// and caches it locally if not already cached.
//
// Parameters:
// - ctx: The context for managing timeouts and cancellation.
// - target: The target platform for which the binary is being downloaded.
//
// Returns:
// - A string representing the path to the cached binary file.
// - An error if any issues occur during the process of fetching or extracting the binary.
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
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, binaryURL, nil)
	if err != nil {
		return "", errors.Errorf("failed to build faucet download request: %w", err)
	}
	resp, err := client.Do(req)
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

// BinaryName generates the binary name by concatenating the faucet binary name and version.
func BinaryName() string {
	return fmt.Sprintf("%s_%s", faucetBinaryName, faucetVersion.String())
}

// binCacheDirPath constructs and returns the path to the binary cache directory based on the configuration directory.
func binCacheDirPath() (string, error) {
	return xfilepath.Join(config.DirPath, xfilepath.Path(binaryCacheDirectory))()
}

// binCachePath constructs the full path to the cached binary and ensures the necessary directories exist.
// Returns the constructed binary path as a string and any error encountered during directory creation.
func binCachePath() (string, error) {
	dirPath, err := binCacheDirPath()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dirPath, BinaryName())
	return path, os.MkdirAll(dirPath, 0o755)
}
