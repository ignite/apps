package tree

import (
	"fmt"
	"strings"
)

// Node represents a node in a tree.
type Node struct {
	Text     string
	Children []*Node
}

// NewNode creates a new node with the given text.
func NewNode(text string) *Node {
	return &Node{
		Text:     text,
		Children: []*Node{},
	}
}

// AddChild adds a child node to the node.
func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

// Format implements fmt.Formatter.
func (n *Node) Format(f fmt.State, _ rune) {
	fprintNode(f, "", n)
}

func fprintNode(f fmt.State, prefix string, n *Node) {
	fmt.Fprintln(f, n.Text)

	width := 2
	if w, ok := f.Width(); ok {
		width = w
	}

	for i, child := range n.Children {
		if child == nil {
			fmt.Fprint(f, prefix)
			if i < len(n.Children)-1 {
				fmt.Fprintln(f, "│")
			}
			continue
		}
		if i < len(n.Children)-1 {
			fmt.Fprintf(f, "%s├%s ", prefix, strings.Repeat("─", width))
			fprintNode(f, prefix+"│"+strings.Repeat(" ", width+1), child)
		} else {
			fmt.Fprintf(f, "%s└%s ", prefix, strings.Repeat("─", width))
			fprintNode(f, prefix+strings.Repeat(" ", width+2), child)
		}
	}
}
