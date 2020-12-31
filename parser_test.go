package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type FakeLexer struct {
	token, cache string
	scanner      *bufio.Scanner
}

func fakeLexer(t *testing.T, dir, filename string) (l Lexer, cleanup func()) {
	tokenSrc := dir + "golden/" + strings.Replace(filename, ".c", ".lex", 1)
	f, err := os.Open(tokenSrc)
	if err != nil {
		t.Fatal(err)
	}
	cleanup = func() { f.Close() }
	l = &FakeLexer{
		scanner: bufio.NewScanner(f),
	}
	return l, cleanup
}

func (l *FakeLexer) Token() string {
	return l.token
}

func (l *FakeLexer) Rewind() {
	l.cache = l.token
}

func (l *FakeLexer) NextToken() string {
	if l.cache != "" {
		defer func() { l.cache = "" }()
		return l.cache
	}
	l.token = ""
	if l.scanner.Scan() {
		l.token = l.scanner.Text()
	}
	return l.token
}

func TestAST(t *testing.T) {
	// read tokens from corresponding */golden/*.lex
	// trying to do something similar for codegen test - go protobuf?

	runTests := func(dir string) {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			t.Run(dir+file.Name(), func(t *testing.T) {
				lexer, cleanup := fakeLexer(t, dir, file.Name())
				defer cleanup()
				//program, err := AST(lexer)
				_ = lexer
			})
		}
	}

	runTests("tests/valid/")
	runTests("tests/invalid/")
}
