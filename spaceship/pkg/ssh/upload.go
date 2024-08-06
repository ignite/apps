package ssh

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// ProgressCallback is a type for the callback function to update the progress.
type ProgressCallback func(uploaded int64, total int64) error

// progressWriter is an io.Writer that calls a progress callback as data is written.
type progressWriter struct {
	bytesUploaded    int64
	totalBytes       int64
	progressCallback ProgressCallback
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.bytesUploaded += int64(n)
	return n, pw.progressCallback(pw.bytesUploaded, pw.totalBytes)
}

// Upload uploads a directory recursively to the remote server with a progress callback.
func (s *SSH) Upload(ctx context.Context, srcPath, dstPath string, progressCallback ProgressCallback) ([]string, error) {
	var (
		totalFiles    int64
		totalBytes    int64
		uploadedBytes int64
	)

	// Count the total number of files and total bytes to be uploaded
	err := filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, err := filepath.Rel(srcPath, path)
			if err != nil {
				return err
			}
			if !strings.HasPrefix(rel, ".") {
				totalFiles++
				totalBytes += info.Size()
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	grp, ctx := errgroup.WithContext(ctx)
	grp.SetLimit(5)

	uploadedFiles := make([]string, 0)
	err = filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			rel, err := filepath.Rel(srcPath, path)
			if err != nil {
				return err
			}
			// skip hidden files and folders.
			if strings.HasPrefix(rel, ".") {
				return nil
			}
			newPath := filepath.Join(dstPath, rel)

			grp.Go(func() error {
				file, err := s.UploadFile(path, newPath, func(bytesUploaded int64, fileTotalBytes int64) error {
					uploadedBytes += bytesUploaded
					// Call the progress callback with the total uploaded bytes
					return progressCallback(uploadedBytes, totalBytes)
				})
				if err != nil {
					return err
				}
				uploadedFiles = append(uploadedFiles, file)
				return nil
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return uploadedFiles, grp.Wait()
}

// UploadFile uploads a single file to the remote server with progress tracking.
func (s *SSH) UploadFile(filePath, dstPath string, progressCallback ProgressCallback) (string, error) {
	dstDir := filepath.Dir(dstPath)
	if err := s.sftpClient.MkdirAll(dstDir); err != nil {
		return "", errors.Wrapf(err, "failed to create destination path %s", dstDir)
	}

	srcPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open source file %s", srcPath)
	}
	defer srcFile.Close()

	fileInfo, err := srcFile.Stat()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get file info for %s", srcPath)
	}
	totalBytes := fileInfo.Size()
	srcReader := io.TeeReader(srcFile, &progressWriter{
		totalBytes:       totalBytes,
		progressCallback: progressCallback,
	})

	dstFile, err := s.sftpClient.Create(dstPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create destination file %s", dstPath)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcReader)
	if err != nil {
		return "", errors.Wrapf(err, "failed to upload file %s to %s", srcPath, dstPath)
	}
	return dstPath, nil
}

// UploadBinary uploads a binary file to the remote server's bin directory
// and sets the appropriate permissions.
func (s *SSH) UploadBinary(srcPath string, progressCallback ProgressCallback) (string, error) {
	var (
		filename = filepath.Base(srcPath)
		binPath  = filepath.Join(s.Bin(), filename)
	)
	if _, err := s.UploadFile(srcPath, binPath, progressCallback); err != nil {
		return "", err
	}

	// give binary permission
	if err := s.sftpClient.Chmod(binPath, 0o755); err != nil {
		return "", err
	}
	return binPath, nil
}

// UploadRunnerScript uploads a runner script to the remote server
// and sets the appropriate permissions.
func (s *SSH) UploadRunnerScript(srcPath string, progressCallback ProgressCallback) (string, error) {
	path := s.RunnerScript()
	if _, err := s.UploadFile(srcPath, s.RunnerScript(), progressCallback); err != nil {
		return "", err
	}

	// give binary permission
	if err := s.sftpClient.Chmod(path, 0o755); err != nil {
		return "", err
	}
	return path, nil
}

// UploadHome uploads the home directory to the remote server.
func (s *SSH) UploadHome(ctx context.Context, srcPath string, progressCallback ProgressCallback) ([]string, error) {
	path := s.Home()
	return s.Upload(ctx, srcPath, path, progressCallback)
}
