package re

import (
	"log"
	"unicode/utf8"
)

func lexString(s string) []token {
	lex := &lexer{input: s, tks: []token{}}
	lex.run()
	return lex.tks
}

type token struct {
	val rune
	tp  tokenType
}

func (t token) String() string {
	return "{'" + string(t.val) + "': " + typePrint[t.tp] + "}"
}

type tokenType int

const (
	char tokenType = iota
	ope
	empty
	end
)

// for printing only
var typePrint = map[tokenType]string{
	char:  "char",
	ope:   "operator",
	empty: "empty string",
	end:   "EOF",
}

var eof = utf8.RuneError

type lexer struct {
	i, end int
	input  string
	tks    []token
	err    error

	lastRuneSize int
}

func (lex *lexer) run() {
	state := any
	for state != nil {
		state = state(lex)
	}
	lex.Emit(eof, end)
}

func (lex *lexer) next() rune {
	if lex.i < len(lex.input) {
		r, size := utf8.DecodeRuneInString(lex.input[lex.end:])
		if r == utf8.RuneError {
			if size == 0 {
				return eof
			}
			if size == 1 {
				log.Fatalf("invalid rune at %v", lex.end)
			}
		}
		lex.end += size
		lex.lastRuneSize = size
		return r
	}
	return eof
}

func (lex *lexer) Unread() {
	lex.end -= lex.lastRuneSize
	lex.lastRuneSize = 0
}

func (lex *lexer) Emit(r rune, tp tokenType) {
	tk := token{val: r, tp: tp}
	lex.tks = append(lex.tks, tk)
}

type lexState func(*lexer) lexState

func any(lex *lexer) lexState {
	r := lex.next()
	for ; r != eof; r = lex.next() {
		switch r {
		case '\n', ' ', '\t':
			continue
		case '\\':
			r := lex.next()
			if r == eof {
				log.Fatal("unexpected EOF in escape character")
			}
			lex.Emit(escape(r))
		case '[':
			lex.Emit(r, ope)
			return insideSet
		case ']':
			log.Fatal("Unclosed brackets")
		case '(', ')', '*', '|':
			lex.Emit(r, ope)
		case eof:
			lex.Emit(r, end)
		default:
			lex.Emit(r, char)
		}
	}
	return nil
}

func insideSet(lex *lexer) lexState {
	r := lex.next()
	if r == '^' {
		lex.Emit(r, ope)
		r = lex.next()
	}
	for ; r != ']'; r = lex.next() {
		switch r {
		case eof:
			log.Fatal("unexpected EOF inside set")
		case '\\':
			r := lex.next()
			if r == eof {
				log.Fatal("unexpected EOF in escape character")
			}
			v, tp := escape(r)
			if tp == empty {
				log.Fatal("Empty string not permitted inside sets. Use [set]|\\e instead.")
			}
			lex.Emit(v, tp)
		case '-':
			lex.Emit(r, ope)
		default:
			lex.Emit(r, char)
		}
	}
	lex.Emit(']', ope)
	return any
}

func escape(r rune) (rune, tokenType) {
	switch r {
	case 's':
		return ' ', char
	case 'n':
		return '\n', char
	case 't':
		return '\t', char
	case 'r':
		return '\r', char
	case 'e', 'E':
		return 0, empty
	default:
		return r, char
	}
}
