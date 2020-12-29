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

	"+": true,
	"-": true,
	"*": true,
	"/": true,
	"!": true,
	"~": true,
}

type Lexer interface {
	NextToken() string
	Rewind()
	Token() string
}

type PlainLexer struct {
	file         *os.File
	token, cache string
}

func NewPlainLexer(filename string) (l *PlainLexer, cleanup func()) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	l = &PlainLexer{file: f}
	return l, func() { f.Close() }
}

func (l *PlainLexer) Rewind() {
	l.cache = l.token
}

func (l *PlainLexer) NextToken() string {
	if len(l.cache) > 0 {
		defer func() { l.cache = "" }()
		return l.cache
	}

	var token strings.Builder
	var err error
	b := make([]byte, 1)
	for err == nil {
		_, err = l.file.Read(b)
		if b[0] == ' ' || b[0] == '\n' || b[0] == '\t' || b[0] == '\r' {
			if token.Len() > 0 {
				l.token = token.String()
				return l.token
			}
		} else if special[string(b)] {
			if token.Len() > 0 {
				l.file.Seek(-1, 1)
				l.token = token.String()
				return l.token
			} else {
				l.token = string(b)
				return l.token
			}
		} else {
			token.WriteByte(b[0])
		}
	}
	l.token = ""
	return ""
}

func (l *PlainLexer) Token() string {
	return l.token
}

// Buffered version
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
