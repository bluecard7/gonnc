package main

import (
	"bytes"
	"fmt"
	"os"
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

	got := make([]string, 0, 4)
	want := []string{"1", "2", "3", "3"}
	fmt.Fprintf(os.Stdin, "1 2 3")
	for i := 0; i < 2; i++ {
		lexer.NextToken()
	}
	lexer.Rewind()
	for lexer.Token() != "" {
		lexer.NextToken()
	}
	if len(got) != len(want) {
		t.Errorf("lengths don't match")
	}
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("Expected %v, got %v", want, got)
		}
	}
}

func TestNextToken(t *testing.T) {
	tests := []struct {
		Program    string
		GoldenFile string
	}{
		{Program: "0", GoldenFile: "return_0"},
		{Program: "5/9*3+3+2*4*4", GoldenFile: ""},
		// spaced out, etc, probably should bring in the tests folder
	}
	for _, test := range tests {
		fmt.Fprintf(os.Stdin, "int main() {%s;}", test.Program)
		tokens := bytes.NewBuffer([]byte{})
		for lexer.NextToken() != "" {
			_, err := got.WriteString(lexer.Token())
			if err != nil {
				t.Fatal(err)
			}
		}
		// if updating, write got to test.GoldenFile
		want := ioutil.ReadFile(test.GoldenFile)
		got := tokens.Bytes()
		if !bytes.Equal(want, got) {
			t.Errorf("Expected:\n%s\nGot:\n%s", want, got)
		}
	}
}
