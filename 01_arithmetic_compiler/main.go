package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var lang_args = []string{}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: calc \"Expr\"")

		fmt.Println(`
Expr := Term {('+' | '-') Term}.
Term ::= Unary {('*' | '/') Unary}
Unary ::= [('+' | '-')] Factor
Factor := '(' Expr ')'
	| number.

number ::= [0-9]+\.?[0-9]*`)
		os.Exit(0)
	}
	str := os.Args[1]
	lang_args = os.Args[2:]
	tks := LexStr(str)
	fmt.Printf("%s\n", tks)
	root := Parse(tks)
	fmt.Println(root)
	fmt.Println(solve(root))
	gen := &CodeGen{}
	b := gen.Generate(root)
	fmt.Println(b)
	m := Machine{}
	m.Run(b)
	alc := &Allocator{
		in:        b,
		Address:   make(map[string]int, 8),
		rbpOffset: 0,
	}
	res := &Resources{
		Available: NewStack(3),
		Next:      make([]int, 3),
		Location:  make(map[string]int, 8),
		Value:     make(map[int]string, 8),
	}
	out := alc.Begin(res)
	fmt.Println(out)
	f, err := os.OpenFile("out.s", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	_, err = f.Write([]byte(out))
	f.Close()
	if err != nil {
		panic(err)
	}
}

func solve(n *node) int {
	switch n.tp {
	case Tnum:
		return StrToFloat(n.val)
	case Targ:
		index := LangArgToIndex(n.val)
		return StrToFloat(lang_args[index])
	}
	if len(n.leafs) == 1 { // unary
		if n.val == "-" {
			return -solve(n.leafs[0])
		}
		return solve(n.leafs[0])
	}
	return DoOp(n.val, solve(n.leafs[0]), solve(n.leafs[1]))
}

func StrToFloat(s string) int {
	out, ok := strconv.ParseInt(s, 10, 64)
	if ok == nil {
		return int(out)
	}
	log.Fatalf("Invalid number: %v", s)
	return 0
}

func DoOp(op string, a, b int) int {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	default:
		log.Fatalf("Invalid operation: %v", op)
		return 0
	}
}

func LangArgToIndex(arg string) int {
	return int((arg[1] - 48) - 1)
}
