package main

import (
	"fmt"
	"strings"
)

type NodeKind int
type NodeValue int64

type ASTNode struct {
	Kind     NodeKind
	Value    NodeValue
	Children []*ASTNode
}

func NewASTNode(kind NodeKind) *ASTNode {
	return &ASTNode{
		Kind:     kind,
		Children: make([]*ASTNode, 0, 2),
	}
}

func (node *ASTNode) AddChildren(children ...*ASTNode) {
	node.Children = append(node.Children, children...)
}

func translateKind(kind NodeKind) string {
	switch kind {
	case PROGRAM:
		return "Program"
	case FUNCTION:
		return "Function"
	case RETURN:
		return "Return"
	case NUM:
		return "number"
	case ADD:
		return "Add"
	case SUB:
		return "Sub"
	case MUL:
		return "Mul"
	case DIV:
		return "Div"
	case NEG:
		return "Negation"

	}
	return fmt.Sprintf("Unknown Kind(%d)", kind)
}

func (node *ASTNode) Print(lvl int) {
	pad := strings.Repeat("\t", lvl)
	fmt.Printf("%sKind: %s\n", pad, translateKind(node.Kind))
	if node.Kind == NUM {
		fmt.Printf("%sValue: %d\n", pad, node.Value)
	}
	for _, child := range node.Children {
		child.Print(lvl + 1)
	}
	fmt.Println()
}
