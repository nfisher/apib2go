package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// completely ripping off Rob Pike's talk. :)

var EOF = rune(0)
var MetaDelim = "---"

func New(filename, input string) *Lexer {
	return &Lexer{
		name:  filename,
		input: input,
		Items: make(chan Item, 2),
	}
}

type ItemType int
type Item struct {
	Type  ItemType
	Value string
}

type Lexer struct {
	name  string
	input string
	start int
	pos   int
	width int
	Items chan Item
}

func (l *Lexer) Emit(t ItemType) {
	l.Items <- Item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) HasPrefix(prefix string) bool {
	return strings.HasPrefix(l.input[l.start:l.pos], prefix)
}

func (l *Lexer) Next() (ch rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}
	ch, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return ch
}

func (l *Lexer) Ignore() {
	l.start = l.pos
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) Peek() rune {
	ch := l.Next()
	l.backup()
	return ch
}

func (l *Lexer) Pos() int {
	return l.pos
}

func (l *Lexer) Accept(valid string) bool {
	if strings.IndexRune(valid, l.Next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) AcceptRun(valid string) {
	for strings.IndexRune(valid, l.Next()) >= 0 {
	}
	l.backup()
}

func (l *Lexer) Run() {
	for fn := LexMetaKey; fn != nil; {
		fn = fn(l)
	}
	close(l.Items)
}

func (l *Lexer) AcceptClasses(cl ...AcceptFn) {
	for {
		r := l.Next()
		if r == EOF {
			break
		}

		accept := false
		for _, c := range cl {
			if c(r) {
				accept = true
				break
			}
		}

		if !accept {
			break
		}
	}
	l.backup()
}

func (l *Lexer) AcceptUntil(stop string) {
	for {
		r := l.Next()
		if r == EOF {
			break
		} else if strings.IndexRune(stop, r) >= 0 {
			break
		}
	}
	l.backup()
}

func (l *Lexer) Errorf(format string, args ...interface{}) StateFn {
	l.Items <- Item{
		ItemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func RuneSet(set string) AcceptFn {
	return func(r rune) bool {
		return (strings.IndexRune(set, r) >= 0)
	}
}

type StateFn func(*Lexer) StateFn
type AcceptFn func(rune) bool
