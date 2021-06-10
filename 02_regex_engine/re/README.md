# RE to DFA

This is a work in progress, it converts regular expressions accepting sets ("[^0-9]", "[a-z]") to DFA using a variation on the subset construction algorithm.

Here's what you can do with sets:

- [^] -> Σ (the entire alphabet accepted by the automaton, here it means the unicode range)
- [^a] -> Σ-a (the entire alphabet excluding 'a')
- [^abc] -> Σ-c-b-c
- [] -> ø (empty set)
- [abc] -> a|b|c
- [a-z] -> a|b|c|...|x|y|z

## Regex Syntax

- alternation: "a|b", "a|b|c"
- concatenation: "ab"
- zero or more: "a\*"
- empty string: "\e"
- grouping: "(a|b)c", "(ab)\*"
- escapes: "\n", "\\\|", "\\\*"
- sets: "[abcdef]", "[0123456789]"
- negated sets: "[\^abc]", "[\^cd]"
- ranges: "[0-9]", "[\^A-Za-z0-9\_]"

Note: it should support unicode, but the literals must be inserted either directly or using go's unicode escapes ("\u00FF" works, \`\u00FF\` doesn't)

## EBNF

Note: I'm using a recursive descent parser instead of the normal reverse polish notation used in Thompson's paper, this is probably dumb, too lazy to change it now.

```ebnf
RE := Expr | ""
Expr := Str {"|" Str}
Str := Rep {Rep}
Rep := Term ["*"]
Term := "(" Expr ")"
	| '[' Set ']'
	| Char
Set := ['^'] { item }
Item := setchar ['-' setchar]


// lexer
SetChar := ['\'] Rune
        | SetRune
Char := ['\'] Rune
	| NormalRune

SetRune	:= [^\\ \-] // all but set operators
NormalRune := [^\\ \| \* \( \) \[ \]] // all but operators
Rune := [\u0000-\uFFFF] // any unicode codepoint
```
