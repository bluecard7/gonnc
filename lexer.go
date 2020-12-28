package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

var special = map[string]bool{
	"{": true,
	"}": true,
	"(": true,
	")": true,
	";": true,
}

type Lexer interface {
	NextToken() string
}

type PlainLexer struct {
	file *os.File
}

func NewPlainLexer(filename string) PlainLexer {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	return PlainLexer{
		file: f,
	}
}

func (l PlainLexer) NextToken() string {
	var token strings.Builder
	var err error
	b := make([]byte, 1)
	for err == nil {
		_, err = l.file.Read(b)
		if b[0] == ' ' || b[0] == '\n' || b[0] == '\t' || b[0] == '\r' {
			if token.Len() > 0 {
				return token.String()
			}
		} else if special[string(b)] {
			if token.Len() > 0 {
				l.file.Seek(-1, 1)
				return token.String()
			} else {
				return string(b)
			}
		} else {
			token.WriteByte(b[0])
		}
	}
	defer l.file.Close()
	return ""
}

type BufLexer struct {
	file   *os.File
	bufRdr *bufio.Reader
}

func NewBufLexer(filename string) BufLexer {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	return BufLexer{
		file:   f,
		bufRdr: bufio.NewReader(f),
	}
}

func (l BufLexer) NextToken() string {
	var (
		token strings.Builder
		r     rune
		err   error
	)
	for err == nil {
		r, _, err = l.bufRdr.ReadRune()
		if unicode.IsSpace(r) {
			if token.Len() > 0 {
				return token.String()
			}
		} else if special[string(r)] {
			if token.Len() > 0 {
				l.bufRdr.UnreadRune()
				return token.String()
			}
			return string(r)
		} else {
			token.WriteRune(r) // check for errors?
		}
	}
	defer l.file.Close()
	if err != io.EOF {
		log.Fatal(err)
	}
	return ""
}
