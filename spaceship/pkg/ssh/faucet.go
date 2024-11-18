package ssh

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

const faucetRepo = "https://github.com/ignite/faucet"

func faucetBinary(target string) (string, error) {
	tempDir, err := os.MkdirTemp("", "faucet")
	if err != nil {
		return "", errors.Errorf("failed to create temp dir: %w", err)
	}

	binaryURL := fmt.Sprintf("%s/releases/latest/download/faucet-%s.tar.gz", faucetRepo, target)
	resp, err := http.Get(binaryURL)
	if err != nil {
		return "", errors.Errorf("failed to download faucet binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("failed to fetch faucet binary, status: %s", resp.Status)
	}

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", errors.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	var binaryPath string

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", errors.Errorf("failed to read tar file: %w", err)
		}

		if header.Typeflag == tar.TypeReg && header.Name == "faucet" {
			binaryPath = filepath.Join(tempDir, header.Name)
			outFile, err := os.Create(binaryPath)
			if err != nil {
				return "", errors.Errorf("failed to create binary file: %w", err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return "", errors.Errorf("failed to write binary file: %w", err)
			}
		}
	}

	if binaryPath == "" {
		return "", errors.Errorf("faucet binary not found in the tar file")
	}

	// Make the binary executable
	if err := os.Chmod(binaryPath, 0o755); err != nil {
		return "", errors.Errorf("failed to make binary executable: %w", err)
	}

	return binaryPath, nil
}
