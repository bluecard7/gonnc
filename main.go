package main

import (
	"log"
	"os"
)

func main() {
	lexer, lexerCleanup := NewPlainLexer(os.Args[1])
	defer lexerCleanup()
	program, err := AST(lexer)
	if err != nil {
		log.Fatal(err)
	}
	//prog.PrintAST(0)
	ASTToASM(program)
}
