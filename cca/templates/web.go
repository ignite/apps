package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:web/*
var web embed.FS

func Write(destinationPath string) error {
	// Walk through the embedded files
	return fs.WalkDir(web, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == "." {
			return nil
		}

		outputPath := filepath.Join(destinationPath, path)

		if d.IsDir() {
			// Create directories as needed
			if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", outputPath, err)
			}
		} else {
			// Write file contents
			fileData, err := web.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read embedded file %s: %w", path, err)
			}

			if err := os.WriteFile(outputPath, fileData, os.ModePerm); err != nil {
				return fmt.Errorf("failed to write file %s: %w", outputPath, err)
			}
		}

		return nil
	})
}
