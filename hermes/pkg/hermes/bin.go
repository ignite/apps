package hermes

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ignite/cli/v28/ignite/config"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/xfilepath"
)

const (
	defaultVersion       = "v1.13.1"
	apiURL               = "https://api.github.com/repos/informalsystems/hermes/releases/tags/"
	binaryCacheDirectory = "apps/hermes/bin"
)

type asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type release struct {
	Assets []asset `json:"assets"`
}

// Maps GOARCH and GOOS to the expected naming used in Hermes release assets.
var archMap = map[string]string{
	"amd64": "x86_64",
	"arm64": "aarch64",
	"386":   "i386",
}

var osMap = map[string]string{
	"darwin": "apple-darwin",
	"linux":  "unknown-linux-gnu",
}

// binCacheDirPath returns the local cache directory path where binaries are stored.
func binCacheDirPath() (string, error) {
	cachePath, err := xfilepath.Join(config.DirPath, xfilepath.Path(binaryCacheDirectory))()
	if err != nil {
		return "", errors.Errorf("failed to construct binary cache directory path: %w", err)
	}
	if err := os.MkdirAll(cachePath, 0o755); err != nil {
		return "", errors.Errorf("failed to create cache directory: %w", err)
	}
	return cachePath, nil
}

// binCachePath returns the full path to the cached Hermes binary for a given version.
func binCachePath(version string) (string, error) {
	cachePath, err := binCacheDirPath()
	if err != nil {
		return "", errors.Errorf("failed to get binary cache directory: %w", err)
	}
	return filepath.Join(cachePath, "hermes-"+version), nil
}

// hermesBin returns the path to the Hermes binary, downloading and extracting it if necessary.
func hermesBin(version string) (string, error) {
	binPath, err := binCachePath(version)
	if err != nil {
		return "", err
	}
	if stat, err := os.Stat(binPath); err == nil && !stat.IsDir() {
		return binPath, nil // already cached
	}
	return fetchBin(version)
}

// fetchBin downloads and extracts the correct Hermes binary for the current platform.
func fetchBin(version string) (string, error) {
	url, err := getHermesAssetURL(version)
	if err != nil {
		return "", err
	}
	return downloadAndExtractHermes(url, version)
}

// getHermesAssetURL queries GitHub Releases API and resolves the download URL for the current system.
func getHermesAssetURL(version string) (string, error) {
	osGo := runtime.GOOS
	archGo := runtime.GOARCH

	archMapped, ok := archMap[archGo]
	if !ok {
		return "", errors.Errorf("unsupported architecture: %s", archGo)
	}

	osMapped, ok := osMap[osGo]
	if !ok {
		return "", errors.Errorf("unsupported OS: %s", osGo)
	}

	expectedAssetName := fmt.Sprintf("hermes-%s-%s-%s.tar.gz", version, archMapped, osMapped)

	// Request release metadata from GitHub
	resp, err := http.Get(apiURL + version)
	if err != nil {
		return "", errors.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", errors.Errorf("failed to decode JSON: %w", err)
	}

	// Find the matching binary for this system
	for _, asset := range rel.Assets {
		if asset.Name == expectedAssetName {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", errors.Errorf("no matching asset found for os=%s arch=%s", osMapped, archMapped)
}

// downloadAndExtractHermes downloads a tar.gz archive and extracts the hermes binary to the cache.
func downloadAndExtractHermes(downloadURL, version string) (string, error) {
	tmpFile := filepath.Join(os.TempDir(), filepath.Base(downloadURL))

	// Download archive
	out, err := os.Create(tmpFile)
	if err != nil {
		return "", errors.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", errors.Errorf("failed to download binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", errors.Errorf("failed to save downloaded binary: %w", err)
	}

	return extractHermesBinary(tmpFile, version)
}

// extractHermesBinary unpacks the hermes binary from the given tar.gz archive.
func extractHermesBinary(tarGzPath, version string) (string, error) {
	f, err := os.Open(tarGzPath)
	if err != nil {
		return "", errors.Errorf("failed to open tar.gz file: %w", err)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return "", errors.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tarReader := tar.NewReader(gzr)
	cachePath, err := binCachePath(version)
	if err != nil {
		return "", err
	}

	// Iterate over files in the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // end of archive
		}
		if err != nil {
			return "", errors.Errorf("error reading tar archive: %w", err)
		}

		// Look for the binary named "hermes"
		if strings.HasSuffix(header.Name, "/hermes") || header.Name == "hermes" {
			outFile, err := os.OpenFile(cachePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
			if err != nil {
				return "", errors.Errorf("failed to create hermes binary: %w", err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return "", errors.Errorf("failed to extract hermes binary: %w", err)
			}
			return cachePath, nil
		}
	}
	return "", errors.Errorf("hermes binary not found in archive")
}

// New returns a usable Hermes instance for the given version (or default version if empty).
func New(version string) (*Hermes, error) {
	if version == "" {
		version = defaultVersion
	}
	binPath, err := hermesBin(version)
	if err != nil {
		return nil, errors.Errorf("failed to get hermes binary: %w", err)
	}
	return &Hermes{path: binPath, version: version}, nil
}
