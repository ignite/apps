package tarball

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

const tarballExt = ".tar.gz"

func ExtractFile(ctx context.Context, file, output string, fileList ...string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ExtractData(ctx, f, output, fileList...)
}

func ExtractData(ctx context.Context, file io.Reader, output string, fileList ...string) ([]string, error) {
	format, reader, err := archiver.Identify(ctx, "", file)
	if err != nil {
		return nil, err
	}
	if format.Extension() != tarballExt {
		return nil, errors.Errorf("unexpected format found: expected=%s actual=%s", tarballExt, format.Extension())
	}

	extracted := make([]string, 0)
	err = format.(archiver.Extractor).Extract(ctx, reader, func(ctx context.Context, f archiver.FileInfo) error {
		fileList := fileList
		if !fileIsIncluded(fileList, f.Name()) {
			return nil
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		newFilePath := filepath.Join(output, f.Name())
		newFile, err := os.Create(newFilePath)
		if err != nil {
			return err
		}
		extracted = append(extracted, newFilePath)
		_, err = newFile.ReadFrom(rc)
		return err
	})
	if err != nil {
		return nil, err
	}
	return extracted, nil
}

// fileIsIncluded returns true if filename is included according to
// filenameList; meaning it is in the list, its parent folder/path
// is in the list, or the list is nil.
func fileIsIncluded(filenameList []string, filename string) bool {
	// include all files if there is no specific list
	if filenameList == nil {
		return true
	}
	for _, fn := range filenameList {
		// exact matches are of course included
		if filename == fn {
			return true
		}
		// also consider the file included if its parent folder/path is in the list
		if strings.HasPrefix(filename, strings.TrimSuffix(fn, "/")+"/") {
			return true
		}
	}
	return false
}
