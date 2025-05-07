package networktypes_test

import (
	"testing"

	launchtypes "github.com/ignite/network/x/launch/types"
	"github.com/stretchr/testify/require"

	"github.com/ignite/apps/network/network/networktypes"
)

func TestVerifyPeerFormat(t *testing.T) {
	tests := []struct {
		name string
		peer launchtypes.Peer
		want bool
	}{
		{
			name: "valid peer connection",
			peer: launchtypes.NewPeerConn("node", "peer:port"),
			want: true,
		},
		{
			name: "peer connection without port",
			peer: launchtypes.NewPeerConn("node", "peer"),
			want: false,
		},
		{
			name: "peer connection without the node address",
			peer: launchtypes.NewPeerConn("node", ":port"),
			want: false,
		},
		{
			name: "peer connection without the separator",
			peer: launchtypes.NewPeerConn("node", "peerport"),
			want: false,
		},
		{
			name: "invalid peer tunnel",
			peer: launchtypes.NewPeerTunnel("", "", ""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := networktypes.VerifyPeerFormat(tt.peer)
			require.Equal(t, tt.want, got)
		})
	}
}
