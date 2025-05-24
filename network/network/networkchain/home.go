package networkchain

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/ignite/apps/network/network/networktypes"
)

// ChainHome returns the default home dir used for a chain from SPN.
func ChainHome(launchID uint64) (path string) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(home, networktypes.SPN, strconv.FormatUint(launchID, 10))
}

// IsChainHomeExist checks if a home with the provided launchID already exist.
func IsChainHomeExist(launchID uint64) (path string, ok bool, err error) {
	home := ChainHome(launchID)

	if _, err := os.Stat(home); os.IsNotExist(err) {
		return home, false, nil
	}

	return home, true, nil
}
