package main

import (
	"log"
	"math"
	"strconv"
)

func solve(n *node) float64 {
	if n.tp == Tnum {
		return StrToFloat(n.val)
	}
	if len(n.kids) == 1 { // unary
		return DoUnary(n.val, solve(n.kids[0]))
	}
	out := solve(n.kids[0])
	for _, kid := range n.kids[1:] {
		out = DoOp(n.val, out, solve(kid))
	}
	return out
}

func StrToFloat(s string) float64 {
	out, ok := strconv.ParseFloat(s, 64)
	if ok == nil {
		return out
	}
	log.Fatalf("Invalid number: %v", s)
	return 0
}

func DoOp(op string, a, b float64) float64 {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	case "%":
		return float64(int(a) % int(b))
	case "^":
		return math.Pow(a, b)
	default:
		log.Fatalf("Invalid operation: %v", op)
		return 0
	}
}

func DoUnary(op string, a float64) float64 {
	switch op {
	case "!":
		return float64(factorial(int(a)))
	}
	log.Fatalf("Invalid operation: %v", op)
	return 0
}

func factorial(a int) int {
	out := 1
	for a > 0 {
		out *= a
		a--
	}
	return out
}
