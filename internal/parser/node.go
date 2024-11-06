package parser

import (
	"fmt"
	"strings"
)

type Node interface {
	Value() string
	Type() string
}

type AstNode struct {
	children []*AstNode
	value    Node
}

type Query *AstNode

func NewAstNode(t Node) *AstNode {
	return &AstNode{
		value:    t,
		children: make([]*AstNode, 0),
	}
}

func (n *AstNode) AppendChild(c *AstNode) {
	n.children = append(n.children, c)
}

func (n *AstNode) String() string {
	return n.stringHelper("", true, true)
}

func (n *AstNode) Value() Node {
	return n.value
}

func (n *AstNode) Children() []*AstNode {
	return n.children
}

func (n *AstNode) stringHelper(prefix string, isLast bool, isRoot bool) string {
	var sb strings.Builder

	sb.WriteString(prefix)
	if isRoot {
	} else if isLast {
		sb.WriteString("└── ")
	} else {
		sb.WriteString("├── ")
	}
	sb.WriteString(fmt.Sprintf("[%s: %s]", n.value.Type(), n.value.Value()))
	sb.WriteString("\n")

	childPrefix := prefix
	if isRoot {
	} else if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "|   "
	}

	for i, child := range n.children {
		sb.WriteString(child.stringHelper(childPrefix, i == len(n.children)-1, false))
	}

	return sb.String()
}
