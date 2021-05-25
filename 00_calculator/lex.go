package main

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

func LexStr(s string) []*lexeme {
	l := &Lexer{
		s:   s,
		tks: make([]*lexeme, 0),
	}
	l.run()
	return l.tks
}

const (
	Tnum lexType = iota
	Tope
	Teof
)

// this represents the eof as a rune. it differs from the lexType above
const eof = utf8.RuneError

var printMap = map[lexType]string{
	Tnum: "int",
	Tope: "ope",
	Teof: "EOF",
}

type lexType int

type lexState func(*Lexer) lexState

type lexeme struct {
	val string
	tp  lexType
}

func (l *lexeme) DebugStr() string {
	return fmt.Sprintf("tk{val: %v, tp: %v}", l.val, printMap[l.tp])
}

func (l *lexeme) String() string {
	if l.tp == Teof {
		return "EOF"
	}
	return l.val
}

type Lexer struct {
	s          string
	start, end int
	tks        []*lexeme

	lastRuneWid int
}

func (l *Lexer) next() rune {
	r, w := utf8.DecodeRuneInString(l.s[l.end:])
	if r == utf8.RuneError && w == 1 {
		log.Fatalf("Invalid UTF8 rune in string. Index: %v", l.end)
	}
	l.end += w
	l.lastRuneWid = w
	return r
}

func (l *Lexer) unread() {
	l.end -= l.lastRuneWid
	l.lastRuneWid = 0
}

func (l *Lexer) ignore() {
	l.start = l.end
}

func (l *Lexer) emit(tp lexType) {
	tk := &lexeme{val: l.s[l.start:l.end], tp: tp}
	l.tks = append(l.tks, tk)
	l.start = l.end
}

func (l *Lexer) accept(s string) bool {
	if strings.ContainsRune(s, l.next()) {
		return true
	}
	l.unread()
	return false
}

func (l *Lexer) acceptRun(s string) {
	for strings.ContainsRune(s, l.next()) {
	}
	l.unread()
}

func (l *Lexer) run() {
	current := any
	for current != nil {
		current = current(l)
	}
}

func any(l *Lexer) lexState {
	r := l.next()
	switch r {
	case ' ', '\n', '\t':
		l.ignore()
		return any
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.unread()
		return number
	case '+', '-', '/', '*', '(', ')', '%', '^', '!':
		l.emit(Tope)
		return any
	case eof:
		l.emit(Teof)
		return nil
	default:
		log.Fatalf("Invalid rune: %v", string(r))
		return nil
	}
}

func number(l *Lexer) lexState {
	l.acceptRun("0123456789")
	if l.accept(".") {
		l.acceptRun("0123456789")
	}
	l.emit(Tnum)
	return any
}
