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

func AST(lexer Lexer) *ASTNode {
	program := NewASTNode("Program")
	for token := lexer.NextToken(); token != ""; token = lexer.NextToken() {
		if token == "int" {
			funcNode := NewASTNode("Function")
			funcNode.Put("Name", lexer.NextToken())
			// parens, func params
			lexer.NextToken()
			lexer.NextToken()
			// {
			lexer.NextToken()
			for token != "}" {
				stmtNode := NewASTNode("Stmt")
				switch lexer.NextToken() {
				case "return":
					stmtNode.Put("StmtType", "Return")
					// would refer to another Node if not a literal value
					stmtNode.Put("Value", lexer.NextToken())
					lexer.NextToken()
				case "}":
					token = "}"
				}
				if token != "}" {
					funcNode.AddChild(stmtNode)
				}
			}
			program.AddChild(funcNode)
		}

	}
	return program
}

func (node *ASTNode) PrintAST(lvl int) {
	for k, v := range node.data {
		fmt.Println(strings.Repeat("\t", lvl), k, ":", v)
	}

	for _, child := range node.children {
		child.PrintAST(lvl + 1)
	}
}
