package main

import (
	"fmt"
	"strconv"
)

const (
	PROGRAM = iota
	FUNCTION
	RETURN
	NUM

	NEG
	NOT
	COMPLEMENT

	ADD
	SUB
	MUL
	DIV
)

func tokenErr(want, got string) error {
	return fmt.Errorf("Expected %s, got %s", want, got)
}

func unknownToken(got string) error {
	return fmt.Errorf("Unknown token: %s", got)
}

func AST(lexer Lexer) (*ASTNode, error) {
	program := NewASTNode(PROGRAM)
	for lexer.NextToken() != "" {
		if lexer.Token() == "int" {
			funcNode := NewASTNode(FUNCTION)
			// IDEA: Map function names to their ASTNodes
			//funcNode.Put("Name", lexer.NextToken())
			lexer.NextToken()
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
				var stmtNode *ASTNode
				switch lexer.Token() {
				case "return":
					stmtNode = NewASTNode(RETURN)
					expr, err := expression(lexer)
					if err != nil {
						return program, err
					}
					stmtNode.AddChildren(expr)
					if lexer.NextToken() != ";" {
						return program, tokenErr(";", lexer.Token())
					}
				default:
					return program, unknownToken(lexer.Token())
				}
				funcNode.AddChildren(stmtNode)
			}
			program.AddChildren(funcNode)
		} else {
			return program, unknownToken(lexer.Token())
		}
	}
	return program, nil
}

func expression(lexer Lexer) (expr *ASTNode, err error) {
	// Placeholder for top level expression parsing
	expr, err = additive(lexer)
	return expr, err
}

func additive(lexer Lexer) (add *ASTNode, err error) {
	mul1, err := multiplicative(lexer)
	if err != nil {
		return mul1, err
	}
	switch lexer.NextToken() {
	case "+":
		add = NewASTNode(ADD)
	case "-":
		add = NewASTNode(SUB)
	default:
		lexer.Rewind()
		return mul1, nil
	}
	mul2, err := multiplicative(lexer)
	if err != nil {
		return mul2, err
	}
	add.AddChildren(mul1, mul2)
	return add, err
}

func multiplicative(lexer Lexer) (mul *ASTNode, err error) {
	prim1, err := primary(lexer)
	if err != nil {
		return prim1, err
	}
	switch lexer.NextToken() {
	case "*":
		mul = NewASTNode(MUL)
	case "/":
		mul = NewASTNode(DIV)
	default:
		lexer.Rewind()
		return prim1, nil
	}
	prim2, err := primary(lexer)
	if err != nil {
		return prim2, err
	}
	mul.AddChildren(prim1, prim2)
	return mul, err
}

func primary(lexer Lexer) (prim *ASTNode, err error) {
	switch lexer.NextToken() {
	case "(":
		prim, err = expression(lexer)
		if lexer.NextToken() != ")" {
			return prim, tokenErr(")", lexer.Token())
		}
	default:
		prim = NewASTNode(NUM)
		v, err := strconv.ParseInt(lexer.Token(), 10, 64)
		if err != nil {
			return prim, err
		}
		prim.Value = NodeValue(v)
	}
	return prim, err
}
