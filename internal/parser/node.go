package parser

import (
	"fmt"
	"strings"
)

type Node interface {
	Value() string
	Type() string
}

type astNode struct {
	children []*astNode
	value    Node
}

type Query *astNode

func NewAstNode(t Node) *astNode {
	return &astNode{
		value:    t,
		children: make([]*astNode, 0),
	}
}

func (n *astNode) AppendChild(c *astNode) {
	n.children = append(n.children, c)
}

func (n *astNode) String() string {
	return n.stringHelper("", true, true)
}

// Recursive helper for the String function
func (n *astNode) stringHelper(prefix string, isLast bool, isRoot bool) string {
	var sb strings.Builder

	// Add the current node's representation
	sb.WriteString(prefix)
	if isRoot {
	} else if isLast {
		sb.WriteString("└── ")
	} else {
		sb.WriteString("├── ")
	}
	sb.WriteString(fmt.Sprintf("[%s: %s]", n.value.Type(), n.value.Value()))
	sb.WriteString("\n")

	// Prepare prefix for the next level
	childPrefix := prefix
	if isRoot {
	} else if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "|   "
	}

	// Recursively print each child with the adjusted prefix
	for i, child := range n.children {
		sb.WriteString(child.stringHelper(childPrefix, i == len(n.children)-1, false))
	}

	return sb.String()
}
