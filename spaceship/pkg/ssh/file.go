package ssh

import (
	"context"
	"fmt"
)

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
