package tarball

import (
	"context"
	"os"
	"path"
	"path/filepath"

	"github.com/mholt/archiver/v4"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

const tarballExt = ".tar.gz"

func Extract(ctx context.Context, file, output string, fileList ...string) ([]string, error) {
	baseName := path.Base(file)
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	format, reader, err := archiver.Identify(baseName, f)
	if err != nil {
		return nil, err
	}
	if format.Name() != tarballExt {
		return nil, errors.Errorf("unexpected format found: expected=%s actual=%s", tarballExt, format.Name())
	}

	extracted := make([]string, 0)
	err = format.(archiver.Extractor).Extract(ctx, reader, fileList, func(_ context.Context, f archiver.File) error {
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
