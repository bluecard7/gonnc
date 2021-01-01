package main

import (
	"encoding/json"
	"fmt"
)

type NodeKind int

var kindReprPairs = []struct {
	kind NodeKind
	name string
}{
	{kind: PROGRAM, name: "Program"},
	{kind: FUNCTION, name: "Function"},
	{kind: RETURN, name: "Return"},
	{kind: NUM, name: "Number"},
	{kind: ADD, name: "Add"},
	{kind: SUB, name: "Subtract"},
	{kind: MUL, name: "Multiply"},
	{kind: DIV, name: "Divide"},
	{kind: NEG, name: "Negate"},
}

func readableKind(kind NodeKind) string {
	for _, pair := range kindReprPairs {
		if pair.kind == kind {
			return pair.name
		}
	}
	return fmt.Sprintf("Unknown Kind(%d)", kind)
}

func compactKind(name string) NodeKind {
	for _, pair := range kindReprPairs {
		if pair.name == name {
			return pair.kind
		}
	}
	return -1 // unknown kind
}

func (k NodeKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(readableKind(k))
}

func (k *NodeKind) UnmarshalJSON(b []byte) error {
	name := string(b)[1 : len(b)-1] // removes "
	*k = compactKind(name)
	return nil
}

type NodeValue int64

type ASTNode struct {
	Kind     NodeKind
	Value    NodeValue  `json:",omitempty"`
	Children []*ASTNode `json:",omitempty"`
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

func (node *ASTNode) JSON() []byte {
	b, err := json.Marshal(node)
	if err != nil {
		panic(err)
	}
	return b
}
