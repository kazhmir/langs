package state

import (
	ir "air/ast"
	T "air/ast/T"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

type Production func(st *State) *ir.Node
type Validator func(*ir.Node) bool

func NewState(input string, filename string) *State {
	st := &State{
		file:  filename,
		input: input,
		line:  1,
		col:   1,
	}
	Next(st)
	return st
}

type State struct {
	Word *ir.Node

	file      string
	line, col int

	start, end   int
	lastRuneSize int
	input        string

	peeked *ir.Node
}

func Consume(st *State) *ir.Node {
	n := st.Word
	Next(st)
	return n
}

func Expect(st *State, tpList ...T.NodeType) *ir.Node {
	for _, tp := range tpList {
		if st.Word.NType == tp {
			return Consume(st)
		}
	}
	fmt.Printf("%v: Expected one of '%v': instead found '%v'\n",
		Place(st),
		T.FmtType(tpList...),
		st.Word.Text)
	os.Exit(1)
	return nil
}

func ExpectProd(st *State, p Production, thing string) *ir.Node {
	n := p(st)
	if n == nil {
		fmt.Printf("%s: Expected %s\n", Place(st), thing)
		os.Exit(1)
	}
	return n
}

func ErrBadSymbol(st *State, symbols ...string) {
	fmt.Printf("%v: Expected one of: '%v', instead found '%v'\n", // TODO: plural/singular
		Place(st),
		strings.Join(symbols, "', '"),
		st.Word.Text,
	)
	os.Exit(1)
}

func Track(st *State, s string) {
	fmt.Printf("%v: %v\n", s, st.Word.String())
}

func Next(l *State) {
	if l.peeked != nil {
		p := l.peeked
		l.peeked = nil
		l.Word = p
		return
	}
	symbol := any(l)
	l.start = l.end
	l.Word = symbol
}

func Peek(l *State) *ir.Node {
	symbol := any(l)
	l.start = l.end
	l.peeked = symbol
	return symbol
}

func Selected(l *State) string {
	return l.input[l.start:l.end]
}

func Place(l *State) string {
	return fmt.Sprintf("%v:%v:%v", l.file, l.line, l.col)
}

func ErrEOF(st *State) {
	fmt.Printf("%s: Unexpected EOF\n", Place(st))
	os.Exit(1)
}

func GetAllTokens(l *State) []*ir.Node {
	out := []*ir.Node{}
	for l.Word.NType != T.EOF {
		out = append(out, l.Word)
		Next(l)
	}
	return out
}

/* RepeatBinary implements the following pattern
for a given Production and Terminal:

	RepeatBinary := Production {Terminal Production}

Validator checks for terminals.
Left to Right precedence
*/
func RepeatBinary(st *State, prod Production, v Validator) *ir.Node {
	last := prod(st)
	for v(st.Word) {
		parent := Consume(st)
		ir.AddLeaf(parent, last)

		newLeaf := prod(st)
		ir.AddLeaf(parent, newLeaf)

		last = parent
	}
	return last
}

/* Repeat implements the following pattern
for a given Production:

	Repeat := {Production}.
*/
func Repeat(st *State, prod Production) []*ir.Node {
	out := []*ir.Node{}
	n := prod(st)
	for n != nil {
		out = append(out, n)
		n = prod(st)
	}
	return out
}

/* RepeatList implements the following pattern
for a given Production and Terminal:

	RepeatBinary := Production {Terminal Production}

Validator checks for terminals.

It differs from RepeatBinary in that it returns a slice
instead of a Tree with precedence
*/
func RepeatList(st *State, prod Production, val Validator) []*ir.Node {
	first := prod(st)
	if first == nil {
		return nil
	}
	out := []*ir.Node{first}
	for val(st.Word) {
		Next(st)
		n := prod(st)
		out = append(out, n)
	}
	return out
}

/*RepeatUnaryRight implements the following pattern
for a given Production:

	RepeatUnaryRight := {Production}.

But returns the first and last item in the tree.

It's Right associative: first->second->last
*/
func RepeatUnaryRight(st *State, f Production) (first, last *ir.Node) {
	first = f(st)
	last = first
	for first != nil {
		n := f(st)
		if n == nil {
			break
		}
		ir.AddLeaf(last, n)
		last = n
	}
	return first, last
}

/*RepeatUnaryLeft implements the following pattern
for a given Production:

	RepeatUnaryLeft := {Production}.

But returns the first and last item in the tree.

It's Left associative: first<-second<-last
*/
func RepeatUnaryLeft(st *State, f Production) (first, last *ir.Node) {
	first = f(st)
	last = first
	for first != nil {
		n := f(st)
		if n == nil {
			break
		}
		ir.AddLeaf(n, last)
		last = n
	}
	return first, last
}

func nextRune(l *State) rune {
	l.col++
	r, size := utf8.DecodeRuneInString(l.input[l.end:])
	if r == utf8.RuneError && size == 1 {
		fmt.Printf("Invalid UTF8 rune in string. Index: %v", l.end)
		os.Exit(1)
	}
	l.end += size
	l.lastRuneSize = size

	return r
}

/*unread decrements the end index by the size of the last rune read,
can only be used once after a Next()*/
func unread(l *State) {
	if l.end > 0 {
		l.end -= l.lastRuneSize
		l.lastRuneSize = 0
		l.col--
	}
}

/*Peek returns the next rune without incrementing the end index*/
func peekRune(l *State) rune {
	r := nextRune(l)
	unread(l)
	return r
}

/*ignore ignores the text previously read*/
func ignore(l *State) {
	l.start = l.end
	l.lastRuneSize = 0
}

func accept(l *State, s string) bool {
	if strings.ContainsRune(s, nextRune(l)) {
		return true
	}
	unread(l)
	return false
}
func acceptRun(l *State, s string) {
	for strings.ContainsRune(s, nextRune(l)) {
	}
	unread(l)
}

func acceptUntil(l *State, s string) {
	for !strings.ContainsRune(s, nextRune(l)) {
	}
	unread(l)
}

const (
	/*eof is equivalent to RuneError, but in this package it only shows up in EoFs
	If the rune is invalid, it panics instead*/
	eof rune = utf8.RuneError
)

var (
	whitespace = `\s\n\t`
	insideStr  = `\"`
	insideRune = `\'`
	digits     = "0123456789"
	letters    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	special    = "_+=<>*"
	idChars   = digits + letters + special
)

func any(st *State) *ir.Node {
	var tp T.NodeType
	var r = nextRune(st)
	for r == ' ' || r == '\n' || r == '\t' {
		if r == '\n' {
			st.line++
			st.col = 0
		}
		ignore(st)
		r = nextRune(st)
	}
	switch r {
	case '0', '1', '2', '3', '4',
		'5', '6', '7', '8', '9':
		return number(st)
	case '"':
		return strLitDouble(st)
	case '\'':
		return strLitSimple(st)
	case '-': // ->
		if nextRune(st) == '>' {
			tp = T.ARROW
		} else {
			unread(st)
			fmt.Printf("%s INVALID SIMBOL: %s\n", Place(st), Selected(st))
			os.Exit(1)
		}
	case ':': // :  ::
		tp = T.COLON
	case '.': // . ..
		tp = T.DOT
	case '{':
		tp = T.LEFTBRACKET
	case '}':
		tp = T.RIGHTBRACKET
	case ',':
		tp = T.COMMA
	case '?':
		tp = T.INTERROGATION
	case eof:
		tp = T.EOF
	default:
		return identifier(st)
	}

	return &ir.Node{
		Text:  Selected(st),
		Col:   st.col,
		Line:  st.line,
		NType: tp,
	}
}

func strLitDouble(st *State) *ir.Node {
	for {
		acceptUntil(st, insideStr)
		r := nextRune(st)
		if r == '"' {
			return &ir.Node{
				Text:  Selected(st),
				NType: T.String,
				Line:  st.line,
				Col:   st.col,
			}
		}
		if r == '\\' {
			nextRune(st) // escaped rune
		}
		unread(st)
	}

}

func strLitSimple(st *State) *ir.Node {
	for {
		acceptUntil(st, insideRune)
		r := nextRune(st)
		if r == '\'' {
			return &ir.Node{
				Text:  Selected(st),
				NType: T.String,
				Line:  st.line,
				Col:   st.col,
			}
		}
		if r == '\\' {
			nextRune(st) // escaped rune
		}
		unread(st)
	}
}

func number(st *State) *ir.Node {
	acceptRun(st, digits)
	if accept(st, ".") {
		acceptRun(st, digits)
	}
	return &ir.Node{
		Text:  Selected(st),
		Line:  st.line,
		Col:   st.col,
		NType: T.Int,
	}
}

func identifier(st *State) *ir.Node {
	acceptRun(st, idChars)
	var tp T.NodeType
	switch Selected(st) {
	case "if":
		tp = T.IF
	case "goto":
		tp = T.GOTO
	case "ret":
		tp = T.RET
	case "true", "false":
		tp = T.Bool
	default:
		tp = T.Id
	}
	return &ir.Node{
		Text:  Selected(st),
		NType: tp,
	}
}
