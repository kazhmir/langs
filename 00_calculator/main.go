package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: calc \"Expr\"")

		fmt.Println(`
Expr ::= Term {("+" | "-") Term}

Term ::= Power {( "*" | "/" | "%" ) Power}

Power ::= Unary { "^" Unary}

Unary ::= Factor ["!"]

Factor ::= "(" Expr ")"
	| SigNum

SigNum ::= [("+" | "-")] Num

 - - - - - - - This is taken care by the lexer

Num ::= {digits} ["." {digits}]

digits ::= '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9'

Note:
Operations only valid in integers ("!", "%") implicitly convert any value to integers.`)
		os.Exit(0)
	}
	str := os.Args[1]
	tks := LexStr(str)
	root := Parse(tks)
	fmt.Println(solve(root))
}
