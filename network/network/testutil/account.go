package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
)

const (
	TestAccountName = "test"
)

// NewTestAccount creates an account for test purposes using in-memory keyring backend.
func NewTestAccount(t *testing.T, name string) cosmosaccount.Account {
	t.Helper()
	r, err := cosmosaccount.NewInMemory()
	assert.NoError(t, err)
	account, _, err := r.Create(name)
	assert.NoError(t, err)
	return account
}
