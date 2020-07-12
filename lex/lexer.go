package lex

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	eof = rune(-1)
)

// Reader ...
func Reader(r io.Reader, opts ...Option) (<-chan Item, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	l := lexer{
		input: unfold(string(b)),
		items: make(chan Item),
	}

	for _, opt := range opts {
		opt(&l)
	}

	go func() {
		for state := lexContentLine; state != nil; {
			state = state(&l)
		}
		close(l.items)
	}()

	return l.items, nil
}

// File ...
func File(filepath string, opts ...Option) (<-chan Item, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Reader(f, opts...)
}

// Text ...
func Text(text string, opts ...Option) (<-chan Item, error) {
	return Reader(strings.NewReader(text))
}

// Bytes ...
func Bytes(b []byte, opts ...Option) (<-chan Item, error) {
	return Reader(bytes.NewReader(b))
}

// Option ...
type Option func(*lexer)

// StrictLineBreaks ...
func StrictLineBreaks(l *lexer) {
	l.strictLineBreaks = true
}

// Must ...
func Must(items <-chan Item, err error) <-chan Item {
	if err != nil {
		panic(fmt.Errorf("lexer: %w", err))
	}
	return items
}

type lexer struct {
	strictLineBreaks bool
	input            string
	start            int
	pos              int
	width            int
	items            chan Item
}

type stateFunc func(*lexer) stateFunc

func (l *lexer) emit(t ItemType) {
	l.items <- Item{
		Type:  t,
		Value: l.input[l.start:l.pos],
	}
	l.start = l.pos
}

func (l *lexer) emitIf(cond bool, t ItemType) {
	if cond {
		l.emit(t)
	}
}

func (l *lexer) emitAdvanced(t ItemType) {
	l.emitIf(l.pos > l.start, t)
}

func (l *lexer) emitEOF() {
	l.ignore()
	l.emit(EOF)
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width

	return
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) hasPrefix(prefix string) bool {
	return strings.HasPrefix(l.input[l.pos:], prefix)
}

func (l *lexer) errorf(format string, args ...interface{}) stateFunc {
	l.items <- Item{
		Type:  Error,
		Value: fmt.Sprintf(format, args...),
	}
	return nil
}

func (l *lexer) unexpected(r rune, valid ...rune) stateFunc {
	svalid := make([]string, len(valid))
	for i, r := range valid {
		svalid[i] = string(r)
	}
	return l.errorf("expected character at pos %d to be one of %v, but got %s", l.pos, svalid, string(r))
}

func (l *lexer) unexpectedEOF() stateFunc {
	return l.errorf("unexpected end of file at pos %d", l.pos)
}

var crlfUnfoldRE = regexp.MustCompile(`\r\n\s`)
var lfUnfoldRE = regexp.MustCompile(`\n\s`)

func unfold(text string) string {
	unfolded := crlfUnfoldRE.ReplaceAllString(text, "")
	return lfUnfoldRE.ReplaceAllString(unfolded, "")
}
