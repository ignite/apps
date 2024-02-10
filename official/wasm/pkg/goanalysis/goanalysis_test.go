package goanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendCode(t *testing.T) {
	type args struct {
		fileContent  string
		functionName string
		codeToInsert string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendCode(tt.args.fileContent, tt.args.functionName, tt.args.codeToInsert)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestAppendImports(t *testing.T) {
	type args struct {
		fileContent      string
		importStatements []string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendImports(tt.args.fileContent, tt.args.importStatements...)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestReplaceCode(t *testing.T) {
	type args struct {
		fileContent     string
		oldFunctionName string
		newFunction     string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceCode(tt.args.fileContent, tt.args.oldFunctionName, tt.args.newFunction)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestReplaceReturn(t *testing.T) {
	type args struct {
		fileContent  string
		functionName string
		returnVars   []string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceReturn(tt.args.fileContent, tt.args.functionName, tt.args.returnVars...)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
