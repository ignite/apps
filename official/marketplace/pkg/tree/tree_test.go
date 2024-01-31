package tree_test

import (
	"fmt"
	"testing"

	"github.com/ignite/apps/official/marketplace/pkg/tree"
	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {
	var (
		assert = assert.New(t)
		node   = tree.NewNode("parent1")
	)
	
	node.AddChild(nil)
	node.AddChild(tree.NewNode("child1"))
	node.AddChild(nil)
	node.AddChild(tree.NewNode("child2"))

	expected := "parent1\n│\n├───── child1\n│\n└───── child2\n"
	assert.Equal(expected, fmt.Sprintf("%5s", node))
}
