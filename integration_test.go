package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

var printAST = flag.Bool("ast", false, "prints AST of program")

// Using tests from https://github.com/nlsandler/write_a_c_compiler
func TestCompiler(t *testing.T) {
	tests := []struct {
		Dir string
	}{
		{Dir: "tests/stage_1/"},
	}
	for _, test := range tests {
		t.Run(test.Dir, func(t *testing.T) {
			testOutput(t, test.Dir, "valid/")
			testOutput(t, test.Dir, "invalid/")
		})
	}
}

func testOutput(t *testing.T, dir, label string) {
	t.Helper()
	path := dir + label
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		filepath := path + file.Name()
		t.Run(label+file.Name(), func(t *testing.T) {
			fmt.Println("[COMPILING]:", filepath)
			lexer, lexerCleanup := NewPlainLexer(filepath)
			defer lexerCleanup()
			program, err := AST(lexer)
			if err != nil {
				log.Println(err)
				return
			}
			if *printAST {
				program.PrintAST(0)
			}
			ASTToASM(program)
		})
	}
}
