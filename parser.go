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

func additive(lexer Lexer) (*ASTNode, error) {
	add, err := multiplicative(lexer)
	for err == nil {
		var kind NodeKind
		switch lexer.NextToken() {
		case "+":
			kind = ADD
		case "-":
			kind = SUB
		default:
			lexer.Rewind()
			return add, err
		}
		mul, err := multiplicative(lexer)
		if err != nil {
			return mul, err
		}
		nextAdd := NewASTNode(kind)
		nextAdd.AddChildren(add, mul)
		add = nextAdd
	}
	return add, err
}

func multiplicative(lexer Lexer) (*ASTNode, error) {
	mul, err := unary(lexer)
	for err == nil {
		var kind NodeKind
		switch lexer.NextToken() {
		case "*":
			kind = MUL
		case "/":
			kind = DIV
		default:
			lexer.Rewind()
			return mul, err
		}
		prim, err := unary(lexer)
		if err != nil {
			return prim, err
		}
		nextMul := NewASTNode(kind)
		nextMul.AddChildren(mul, prim)
		mul = nextMul
	}
	return mul, err
}

func unary(lexer Lexer) (u *ASTNode, err error) {
	switch lexer.NextToken() {
	case "+":
		u, err = unary(lexer)
	case "-":
		operand, err := unary(lexer)
		if err != nil {
			return operand, err
		}
		u = NewASTNode(NEG)
		u.AddChildren(operand)
	default:
		lexer.Rewind()
		u, err = primary(lexer)
	}
	return u, err
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
