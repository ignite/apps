package tx

import (
	"context"
	"errors"

	"github.com/spf13/pflag"

	"cosmossdk.io/core/address"
	"github.com/ignite/apps/connect/internal/autocli/keyring"

	"github.com/cosmos/cosmos-sdk/codec"
)

// ContextKey is a key used to store and retrieve Context from a Go context.Context.
var ContextKey contextKey

// contextKey is an empty struct used as a key type for storing Context in a context.Context.
type contextKey struct{}

// Context represents the client context used in autocli commands.
// It contains various components needed for command execution.
type Context struct {
	Flags *pflag.FlagSet

	AddressCodec          address.Codec
	ValidatorAddressCodec address.Codec
	ConsensusAddressCodec address.Codec

	Cdc     codec.Codec
	Keyring keyring.Keyring
}

// SetContext sets client context in the go context.
func SetContext(goCtx context.Context, ctx Context) context.Context {
	return context.WithValue(goCtx, ContextKey, ctx)
}

// GetContext gets the context from the go context.
func GetContext(goCtx context.Context) (Context, error) {
	if c := goCtx.Value(ContextKey); c != nil {
		ctx, ok := c.(Context)
		if !ok {
			return Context{}, errors.New("context value is not of type autocli context")
		}

		return ctx, nil
	}

	return Context{}, errors.New("context does not contain autocli context value")
}
