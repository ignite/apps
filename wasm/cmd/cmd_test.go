package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_relativePath(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)
	relativeRoot, err := filepath.Rel(pwd, "/")
	require.NoError(t, err)

	tests := []struct {
		name    string
		appPath string
		want    string
	}{
		{
			name:    "Relative path within current directory",
			appPath: "subdir/file.txt",
			want:    "subdir/file.txt",
		},
		{
			name:    "Relative path outside current directory",
			appPath: "/path/file.txt",
			want:    filepath.Join(relativeRoot, "path/file.txt"),
		},
		{
			name:    "App path is current directory",
			appPath: ".",
			want:    ".",
		},
		{
			name:    "App path is parent directory",
			appPath: "..",
			want:    "..",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			absPath, err := filepath.Abs(tt.appPath)
			require.NoError(t, err)
			got, err := relativePath(absPath)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
