package tree_test

import (
	"fmt"
	"testing"

	"github.com/ignite/apps/official/marketplace/pkg/tree"
	"github.com/stretchr/testify/require"
)

func TestTree(t *testing.T) {
	var (
		require = require.New(t)
		node    = tree.NewNode("parent1")
	)

	node.AddChild(nil)
	node.AddChild(tree.NewNode("child1"))
	node.AddChild(nil)
	node.AddChild(tree.NewNode("child2"))

	expected := "parent1\n│\n├───── child1\n│\n└───── child2\n"
	require.Equal(expected, fmt.Sprintf("%5s", node))
}
