package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
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

var (
	viewAST   = flag.Bool("v-ast", false, "prints generated AST")
	updateAST = flag.Bool("u-ast", false, "updates json files containing AST")
)

func TestAST(t *testing.T) {
	// generates dst to write AST json representation based on name of file being tested
	jsonFilepath := func(dir, filename string) string {
		return dir + "golden/" + strings.Replace(filename, ".c", ".ast", 1)
	}
	// trying to do something similar for codegen test
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
				program, _ := AST(lexer)
				if *viewAST {
					fmt.Println(string(program.JSON()))
					return
				}
				goldenFile := jsonFilepath(dir, file.Name())
				if *updateAST {
					ioutil.WriteFile(goldenFile, program.JSON(), 0644)
					return
				}
				want, err := ioutil.ReadFile(goldenFile)
				if err != nil {
					t.Fatal(err)
				}
				if got := program.JSON(); !bytes.Equal(want, got) {
					t.Errorf("Expected:\n%s\nGot\n%s\n", want, got)
				}
			})
		}
	}
	runTests("tests/valid/")
	runTests("tests/invalid/")
}
