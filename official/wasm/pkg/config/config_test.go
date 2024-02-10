package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddWasm(t *testing.T) {
	type args struct {
		configPath string
		options    []Option
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AddWasm(tt.args.configPath, tt.args.options...)
			require.Equal(t, tt.err, err)
		})
	}
}

func Test_hasWasm(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasWasm(tt.args.configPath); got != tt.want {
				t.Errorf("hasWasm() = %v, want %v", got, tt.want)
			}
		})
	}
}
