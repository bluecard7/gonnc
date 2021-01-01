package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestToken(t *testing.T) {
	want := "TOKEN"
	lexer := &BufLexer{token: want}
	if got := lexer.Token(); want != got {
		t.Errorf("Expected %s, got %s", want, got)
	}
}

func TestRewind(t *testing.T) {
	// currently simulates a read, could just instantiate with bufio.Reader around String reader
	want := "TOKEN"
	lexer := &BufLexer{token: want}
	lexer.Rewind()
	if got := lexer.NextToken(); want != got {
		t.Errorf("Expected %s, got %s", want, got)
	}
}

var (
	viewLex   = flag.Bool("v-lex", false, "view lexed tokens")
	updateLex = flag.Bool("u-lex", false, "update expected tokens")
)

func TestNextToken(t *testing.T) {
	lexResultPath := func(dir, filename string) string {
		return dir + "golden/" + strings.Replace(filename, ".c", ".lex", 1)
	}
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
				lexer, cleanup := NewLexer(dir + file.Name())
				defer cleanup()
				tokens := bytes.NewBuffer(make([]byte, 0, 4096))
				goldenFilepath := lexResultPath(dir, file.Name())

				var goldenFile *os.File
				if *updateLex {
					goldenFile, err = os.OpenFile(goldenFilepath, os.O_RDWR|os.O_CREATE, 0755)
					if err != nil {
						t.Fatal(err)
					}
				}
				for lexer.NextToken() != "" {
					if *viewLex {
						fmt.Println(lexer.Token())
						continue
					}
					if *updateLex {
						goldenFile.WriteString(lexer.Token() + "\n")
						continue
					}
					_, err := tokens.WriteString(lexer.Token() + "\n")
					if err != nil {
						t.Fatal(err)
					}
				}
				if *viewLex || *updateLex {
					return
				}
				want, err := ioutil.ReadFile(goldenFilepath)
				if err != nil {
					t.Fatal(err)
				}
				if got := tokens.Bytes(); !bytes.Equal(want, got) {
					t.Errorf("Expected:\n%s\nGot:\n%s\n", want, got)
				}
			})

		}
	}
	runTests("tests/valid/")
	runTests("tests/invalid/")
}
