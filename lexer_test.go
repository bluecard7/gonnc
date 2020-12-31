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
	lexer, cleanup := NewLexer("stdin")
	defer cleanup()

	var got []string
	want := []string{"1", "2", "3", "3"}
	// doesn't look like os.Stdin can be written to and read in this manner
	fmt.Fprintf(os.Stdin, "1 2 3")

	for i := 0; i < 2; i++ {
		got = append(got, lexer.NextToken())
	}
	lexer.Rewind()
	for lexer.Token() != "" {
		got = append(got, lexer.NextToken())
	}
	if len(got) != len(want) {
		t.Fatalf("Expected %v, got %v", want, got)
	}
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("Expected %v, got %v", want, got)
		}
	}
}

var update = flag.Bool("u", false, "update goldenfiles")

func TestNextToken(t *testing.T) {
	goldenFilePath := func(dir, filename string) string {
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
				for lexer.NextToken() != "" {
					_, err := tokens.WriteString(lexer.Token() + "\n")
					if err != nil {
						t.Fatal(err)
					}
				}
				gold := goldenFilePath(dir, file.Name())
				if *update {
					ioutil.WriteFile(gold, tokens.Bytes(), 0644)
					return
				}
				want, err := ioutil.ReadFile(gold)
				if err != nil {
					t.Fatal(err)
				}
				got := tokens.Bytes()
				if !bytes.Equal(want, got) {
					t.Errorf("Expected:\n%s\nGot:\n%s\n", want, got)
				}
			})

		}
	}
	runTests("tests/valid/")
	runTests("tests/invalid/")
}
