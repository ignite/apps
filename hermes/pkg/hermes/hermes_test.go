package hermes

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalResult(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		v    any
		want interface{}
		err  error
	}{
		{
			name: "valid result",
			data: []byte(`{"status": "success", "result": {"wallet":{"account":"cosmos139asl6de8mzxedvvxatp2wdna2n6vy3af62srg","address_type":"Cosmos"}}}`),
			v:    &KeysListResult{},
			want: &KeysListResult{
				Wallet{
					Account:     "cosmos139asl6de8mzxedvvxatp2wdna2n6vy3af62srg",
					AddressType: "Cosmos",
				},
			},
		},
		{
			name: "invalid unmarshall object",
			data: []byte(`{"status": "success", "result": {"wallet":{"account":"cosmos139asl6de8mzxedvvxatp2wdna2n6vy3af62srg","address_type":"Cosmos"}}}`),
			v:    &ClientResult{},
			want: &ClientResult{},
		},
		{
			name: "error result",
			data: []byte(`{"status": "error", "result": {"wallet":{"account":"cosmos139asl6de8mzxedvvxatp2wdna2n6vy3af62srg","address_type":"Cosmos"}}}`),
			v:    &KeysListResult{},
			err:  errors.New(`error result (*hermes.KeysListResult) error: {"wallet":{"account":"cosmos139asl6de8mzxedvvxatp2wdna2n6vy3af62srg","address_type":"Cosmos"}}`),
		},
		{
			name: "unmarshal error",
			data: []byte(`{"status": "success", "result": {"wallet":{"account":"cosmos139asl6de8mzxedvvxatp2wdna2n6vy3af62srg","address_type":"Cosmos"}}}`),
			v:    "",
			err:  &json.InvalidUnmarshalError{Type: reflect.TypeOf("")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalResult(tt.data, tt.v)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.want, tt.v)
		})
	}
}

func TestValidateResult(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "valid result",
			data: []byte(`{"status": "success", "result": "some data"}`),
		},
		{
			name: "error result",
			data: []byte(`{"status": "error", "result": "error data"}`),
			err:  errors.New(`result error: "error data"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateResult(tt.data)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
