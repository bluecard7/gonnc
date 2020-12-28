package main

import (
	"fmt"
	"strings"
)

type ASTNode struct {
	data     map[string]string
	children []*ASTNode
}

func NewASTNode(nodeType string) *ASTNode {
	node := &ASTNode{
		data:     make(map[string]string),
		children: []*ASTNode{},
	}
	node.Put("Type", nodeType)
	return node
}

func (node *ASTNode) Put(key, val string) {
	node.data[key] = val
}

func (node *ASTNode) Get(key string) string {
	val, ok := node.data[key]
	if !ok {
		return ""
	}
	return val
}

func (node *ASTNode) AddChild(child *ASTNode) {
	node.children = append(node.children, child)
}

func tokenErr(want, got string) error {
	return fmt.Errorf("Expected %s, got %s", want, got)
}

func AST(lexer Lexer) (*ASTNode, error) {
	program := NewASTNode("Program")
	for lexer.NextToken() != "" {
		if lexer.Token() == "int" {
			funcNode := NewASTNode("Function")
			funcNode.Put("Name", lexer.NextToken())
			// (...params)
			if lexer.NextToken() != "(" {
				return program, tokenErr("(", lexer.Token())
			}
			if lexer.NextToken() != ")" {
				return program, tokenErr(")", lexer.Token())
			}
			// func definition
			if lexer.NextToken() != "{" {
				return program, tokenErr("{", lexer.Token())
			}
			for lexer.NextToken() != "}" {
				stmtNode := NewASTNode("Stmt")
				switch lexer.Token() {
				case "return":
					stmtNode.Put("StmtType", "Return")
					// TODO:: recursively descend if not a literal value
					stmtNode.Put("Value", lexer.NextToken())
					if lexer.NextToken() != ";" {
						return program, tokenErr(";", lexer.Token())
					}
				}
				funcNode.AddChild(stmtNode)
			}
			program.AddChild(funcNode)
		} else {
			return program, fmt.Errorf("Unidentified token: %s", lexer.Token())
		}
	}
	return program, nil
}

func (node *ASTNode) PrintAST(lvl int) {
	for k, v := range node.data {
		fmt.Println(strings.Repeat("\t", lvl), k, ":", v)
	}

	for _, child := range node.children {
		child.PrintAST(lvl + 1)
	}
}
