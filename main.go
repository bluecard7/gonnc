package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	srcFilename = flag.String("src", "", "program file")
	dstFilename = flag.String("dst", "", "file to write assembly to")
	upToLex     = flag.Bool("l", false, "stop after lexing input, prints tokens")
	upToParse   = flag.Bool("p", false, "stop after parsing tokens")
	printAST    = flag.Bool("ast", false, "print AST")
)

func main() {
	flag.Parse()
	lexer, cleanup := NewLexer(*srcFilename)
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

	dst := os.Stdout
	if *dstFilename != "" {
		dst, err = os.OpenFile(*dstFilename, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	Compile(dst, program)
}
