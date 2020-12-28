package main

import (
	"os"
)

func main() {
	lexer := NewPlainLexer(os.Args[1])
	/*
		l := NewBufLexer(os.Args[1])
		for token := l.NextToken(); token != ""; token = l.NextToken() {
			fmt.Println(token)
		}
	*/
	prog := AST(lexer)
	//prog.PrintAST(0)
	ASTToASM(prog)
}
