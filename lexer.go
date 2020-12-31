package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

type Lexer interface {
	NextToken() string
	Rewind()
	Token() string
}

type BufLexer struct {
	bufRdr       *bufio.Reader
	token, cache string
}

func NewLexer(filename string) (l Lexer, cleanup func()) {
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
	l = &BufLexer{bufRdr: bufio.NewReader(f)}
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
		// TODO:: run profiler and see how many allocation + how to reduce
		sb  strings.Builder
		r   rune
		err error
	)
	for err == nil {
		r, _, err = l.bufRdr.ReadRune()
		if err != nil {
			break
		}
		if unicode.IsSpace(r) {
			if sb.Len() > 0 {
				l.token = sb.String()
				return l.token
			}
			continue
		}

		/* How to do !=, ==, <=, <?
		var buf []byte
		buf, err = l.bufRdr.Peek(utf8.UTFMax)
		d, size = utf8.DecodeRune(buf)
		then check string([]rune{r, d}) == "<=", "!=", etc
		{
		*/
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			l.token = string(r)
			// can I combine these cases? what I'm really checking is if
			// the token is finished
			if sb.Len() > 0 {
				l.bufRdr.UnreadRune()
				l.token = sb.String()
			}
			return l.token
		}
		sb.WriteRune(r)
	}
	if err != io.EOF {
		log.Fatal(err)
	}
	l.token = ""
	return ""
}
