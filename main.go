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
	lexer, cleanup := NewLexer(os.Args[len(os.Args)-1])
	defer cleanup()
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
		b := program.JSON()
		fmt.Println(string(b))
	}
	if *upToParse {
		return
	}
	Compile(program)
}
