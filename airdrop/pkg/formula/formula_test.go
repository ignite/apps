package formula

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestValue_Calculate(t *testing.T) {
	type args struct {
		amount math.Int
		staked math.Int
	}
	tests := []struct {
		name  string
		value Value
		args  args
		want  math.Int
	}{
		{
			name: "valid formula 1",
			value: Value{
				Type:   Quadratic,
				Value:  110,
				Ignore: 1,
			},
			args: args{
				amount: math.NewInt(330),
				staked: math.NewInt(30),
			},
			want: math.NewInt(198),
		},
		{
			name: "valid formula 2",
			value: Value{
				Type:   Quadratic,
				Value:  2,
				Ignore: 1,
			},
			args: args{
				amount: math.NewInt(330),
				staked: math.NewInt(200),
			},
			want: math.NewInt(54),
		},
		{
			name: "valid formula without stake",
			value: Value{
				Type:   Quadratic,
				Value:  19,
				Ignore: 1,
			},
			args: args{
				amount: math.NewInt(223),
				staked: math.NewInt(0),
			},
			want: math.NewInt(14),
		},
		{
			name: "zero values",
			value: Value{
				Type:   Quadratic,
				Value:  0,
				Ignore: 0,
			},
			args: args{
				amount: math.NewInt(0),
				staked: math.NewInt(0),
			},
			want: math.NewInt(0),
		},
		{
			name: "invalid formula",
			value: Value{
				Type:   "invalid",
				Value:  110,
				Ignore: 2,
			},
			args: args{
				amount: math.NewInt(330),
				staked: math.NewInt(30),
			},
			want: math.NewInt(0),
		},
		{
			name: "zero amount",
			value: Value{
				Type:   "invalid",
				Value:  110,
				Ignore: 2,
			},
			args: args{
				amount: math.NewInt(0),
				staked: math.NewInt(30),
			},
			want: math.NewInt(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.value.Calculate(tt.args.amount, tt.args.staked)
			require.Equal(t, tt.want, got)
		})
	}
}
