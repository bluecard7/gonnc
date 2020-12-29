package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	upToLex   = flag.Bool("l", false, "stop after lexing input")
	upToParse = flag.Bool("p", false, "stop after parsing tokens")
	printAST  = flag.Bool("ast", false, "print AST")
)

func main() {
	flag.Parse()
	lexer, lexerCleanup := NewPlainLexer(os.Args[len(os.Args)-1])
	defer lexerCleanup()
	if *upToLex {
		for lexer.NextToken() != "" {
			fmt.Println(lexer.Token())
		}
		return
	}
	program, err := AST(lexer)
	if err != nil {
		log.Fatal(err)
	}
	if *upToParse || *printAST {
		program.Print(0)
	}
	if *upToParse {
		return
	}
	Compile(program)
}
