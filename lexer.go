package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"
)

type specialtokens []string

func punctuators() specialtokens {
	var tokens = []string{"{", "}", "(", ")", ";", "+", "-", "*", "/", "!", "~"}
	sort.Strings(tokens)
	return tokens
}

func (s specialtokens) has(target string) bool {
	i := sort.SearchStrings(s, target)
	return i < len(s) && s[i] == target
}

type Lexer interface {
	NextToken() string
	Rewind()
	Token() string
}

type BufLexer struct {
	bufRdr       *bufio.Reader
	punctuators  specialtokens
	token, cache string
}

func NewBufLexer(filename string) (l *BufLexer, cleanup func()) {
	f := os.Stdin
	cleanup = func() {}
	if filename != "stdin" {
		var err error
		f, err = os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		cleanup = func() { f.Close() }
	}
	l = &BufLexer{
		bufRdr:      bufio.NewReader(f),
		punctuators: punctuators(),
	}
	return l, cleanup
}

func (l *BufLexer) Token() string {
	return l.token
}

func (l *BufLexer) Rewind() {
	l.cache = l.token
}

func (l *BufLexer) NextToken() string {
	if len(l.cache) > 0 {
		defer func() { l.cache = "" }()
		return l.cache
	}
	var (
		token strings.Builder
		r     rune
		err   error
	)
	for err == nil {
		r, _, err = l.bufRdr.ReadRune()
		if err != nil {
			break
		}
		if unicode.IsSpace(r) {
			if token.Len() > 0 {
				l.token = token.String()
				return l.token
			}
		} else if l.punctuators.has(string(r)) {
			l.token = string(r)
			if token.Len() > 0 {
				l.bufRdr.UnreadRune()
				l.token = token.String()
			}
			return l.token
		} else {
			token.WriteRune(r) // check for errors?
		}
	}
	if err != io.EOF {
		log.Fatal(err)
	}
	l.token = ""
	return ""
}
