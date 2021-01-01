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

// not exactly, want to rename to something like pathToOutputFromStage
func pathToGoldenfile(dir, filename, fileExt string) string {
	// then based on stage, select proper fileExt
	return dir + "golden/" + strings.ReplaceAll(filename, ".c", fileExt)
}

// just using in lexer test for now

func runTestsOnFiles(t *testing.T, dir string, testrun func(t *testing.T, dir, filename string)) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dir)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		t.Log(file.Name())
		testrun(t, dir, file.Name())
	}
}

func TestNextToken(t *testing.T) {
	testrun := func(t *testing.T, dir, filename string) {
		lexer, cleanup := NewLexer(dir + filename)
		defer cleanup()
		tokens := bytes.NewBuffer(make([]byte, 0, 4096))
		goldenfilePath := pathToGoldenfile(dir, filename, ".lex")

		// just revert back to old way, just write all bytes at once
		var goldenfile *os.File
		var err error
		if *updateLex {
			goldenfile, err = os.OpenFile(goldenfilePath, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				t.Fatal(err)
			}
			defer goldenfile.Close()
		}
		for lexer.NextToken() != "" {
			if *viewLex {
				fmt.Println(lexer.Token())
				continue
			}
			// could combine tokens and goldenFile to a io.Writer...
			if *updateLex {
				goldenfile.WriteString(lexer.Token() + "\n")
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
		want, err := ioutil.ReadFile(goldenfilePath)
		if err != nil {
			t.Fatal(err)
		}
		if got := tokens.Bytes(); !bytes.Equal(want, got) {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", want, got)
		}
	}
	runTestsOnFiles(t, "tests/valid/", testrun)
	runTestsOnFiles(t, "tests/invalid/", testrun)
	/*
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
					goldenfilePath := pathToGoldenfile(dir, file.Name(), ".lex")

					var goldenfile *os.File
					if *updateLex {
						goldenfile, err = os.OpenFile(goldenfilePath, os.O_RDWR|os.O_CREATE, 0755)
						if err != nil {
							t.Fatal(err)
						}
					}
					for lexer.NextToken() != "" {
						if *viewLex {
							fmt.Println(lexer.Token())
							continue
						}
						// could combine tokens and goldenFile to a io.Writer...
						if *updateLex {
							goldenfile.WriteString(lexer.Token() + "\n")
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
					want, err := ioutil.ReadFile(goldenfilePath)
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
	*/
}
