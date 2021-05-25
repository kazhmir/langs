# Action Triggering Automaton

This package was made for educational purposes, it shouldn't be used in production. It uses a DFA based regular expression engine, that does not backtrack, and is guaranteed to read the input in linear time. It's primary goal is to experiment with the concept of ```pattern -> action```, used in some tools like AWK.

The inner workings are simple, it appends an action to the accepting state of every automaton, as the Run method goes through the automaton, these actions are activated. The automaton is built by first parsing the pattern into an Abstract Syntax Tree, then with this tree, a NFA is built using Thompson's Construction. This NFA is then converted into DFA by a variation of the Powerset Algorithm, beware that this conversion can take exponential time.

## Installation

To use you can simply ```go get``` the package:

```go get github.com/kazhmir/ATA```

## Regex Syntax

The Regex engine supports: 
- alternation: "a|b", "a|b|c"
- concatenation: "ab"
- zero or more: "a\*"
- empty string: "\e"
- empty set: "", "[]"
- grouping: "(a|b)c", "(ab)\*"
- escapes: "\t", "\n", "\\\|" 
- convenience groups: "\n", "\N", "\w", "\W", "\s", "\S", "\d", "\D"
- dot: '.' ('[\^\n\r]')
- sets: "[abcdef]", "[0123456789]"
- negated sets: "[\^abc]", "[\^cd]"
- ranges: "[0-9]", "[\^A-Za-z0-9\_]"
- unicode: "😉", "π|τ"

Will be added in the future: 
- 0 or 1: "a?" == "a|\e"
- 1 or more: "a+" == "aa*"

Notes:
	- Whenever the pattern can match an empty string, it will also match EOF.
	- It matches ALL SUBMATCHES (not only the longest and smallest). This means the pattern "ab|bc" running on the string "abc" will match {"ab", "bc"}. I still didn't figure out a way to find the longest match linearly, this is the only thing impeding this package to be able to generate lexers.
	- Performance profile is quite different from the standard regexp package, orders of magnitude slower, but still linear.

## EBNF

The terminals are either *sets* or operators. The parser sees sets "[a]" and chars "a" as the same thing.

```ebnf
RE 	:= Expr | ""
Expr 	:= Str {"|" Str}
Str 	:= Rep {Rep}
Rep 	:= Term ["*"]
Term 	:= "(" Expr ")"
	| '[' Set ']'
	| Char
Set   	:= ['^'] { item }
Item 	:= setchar ['-' setchar]


// lexer
SetChar := ['\'] Rune
        | SetRune
Char 	:= ['\'] Rune
	| NormalRune
SetRune  	:= [^\\ \-] // all but set operators
NormalRune 	:= [^\\ \| \* \( \) \[ \]] // all but operators
Rune 		:= [\u0000-\uFFFF] // any unicode codepoint
```

The Lexer also has the responsability of breaking down convenience groups into sets ("." --> "[^\n\r]", "\w" --> "[a-zA-Z_]"...), escaping runes ("\|", "\\" etc) and expanding ranges ("[a-e]" --> "[abcde]").

## Usage

To use the package, you have to compile the ```Machine```. You can do this with one of the two functions:

```go
func Build(map[string]Action) *Machine
func BuildOne(string, Action) *Machine
```

A Machine is an automaton, it simply holds the starting state, the pattern used to create it and the structure:

```go
type Machine struct {
	Start *state
	Pattern string
	Syntax map[string]Action
}
```

The Action is simply a function of type ```func(*Match) (stop bool)```.

```go
func (*Machine) Run(input io.RuneScanner) error
func (*Machine) Debug(input string) error
```

Two convenience functions are also provided, but you can write others easily, here's ```FindAllString```

```go
func FindAllString(pattern string, txt io.RuneReader) ([]string, error) {
	out := make([]string, 0)
	act := func(mat *Match) bool {
		out = append(out, mat.S)
		return false
	}
	m := BuildOne(pattern, act)
	err := m.Run(txt)
	return out, err
}
```

Further examples:
- [Benchmark against std regexp](./examples/ata.go)
- [grep](./examples/grep.go)
- [Concurrent matching of multiple files](./examples/search.go)
- [Others](./examples)
