package main

import (
	"fmt"
	"os"
)

type node struct {
	*lexeme
	kids []*node
}

func newNode(l *lexeme) *node {
	return &node{
		lexeme: l,
		kids:   []*node{},
	}
}

func (n *node) newLeaf(kid *node) *node {
	n.kids = append(n.kids, kid)
	return kid
}

func (n *node) String() string {
	return n.ast(0)
}

func (n *node) ast(i int) string {
	output := n.lexeme.String() + "\n"
	for _, kid := range n.kids {
		output += indent(i) + kid.ast(i+1)
	}
	return output
}

func indent(n int) string {
	output := ""
	for i := -1; i < n-1; i++ {
		output += "   "
	}
	output += "└─>"
	return output
}

func Parse(tks []*lexeme) *node {
	p := &Parser{
		tks:  tks,
		word: tks[0],
	}
	return p.Expr()
}

type Parser struct {
	i    int
	tks  []*lexeme
	word *lexeme
}

func (p *Parser) next() {
	if p.i < len(p.tks)-1 {
		p.i++
		p.word = p.tks[p.i]
	}
}

func (p *Parser) previous() {
	if p.i > 0 {
		p.i--
		p.word = p.tks[p.i]
	}
}

func (p *Parser) Fail() *node {
	fmt.Printf("Invalid Syntax at: '%v'. Token index %v\n", p.word, p.i)
	os.Exit(0)
	return nil
}

/*Whenever we sucessfully match a terminal p.Next will be present in the same block
 */
func (p *Parser) Expr() *node {
	last := p.Term()
	for p.word.val == "+" || p.word.val == "-" {
		parent := newNode(p.word)
		parent.newLeaf(last)
		p.next()
		parent.newLeaf(p.Term())
		last = parent
	}
	if p.word.val == ")" || p.word.tp == Teof {
		p.next()
		return last
	}
	p.Fail()
	return nil
}

func (p *Parser) Term() *node {
	last := p.Power()
	for p.word.val == "*" || p.word.val == "/" || p.word.val == "%" {
		parent := newNode(p.word)
		parent.newLeaf(last)
		p.next()
		parent.newLeaf(p.Power())
		last = parent
	}
	return last
}

func (p *Parser) Power() *node {
	last := p.Unary()
	if p.word.val == "^" {
		parent := newNode(p.word) // "^" node
		parent.newLeaf(last)
		p.next()
		parent.newLeaf(p.Power())
		return parent
	}
	return last
}

func (p *Parser) Unary() *node {
	last := p.Factor()
	if p.word.val == "!" {
		parent := newNode(p.word)
		p.next()
		parent.newLeaf(last)
		return parent
	}
	return last
}

func (p *Parser) Factor() *node {
	if p.word.val == "(" {
		p.next()
		return p.Expr()
	}
	return p.Num()
}

func (p *Parser) Num() *node {
	sig := "" // optional signal
	if p.word.val == "+" || p.word.val == "-" {
		sig = p.word.val
		p.next()
	}
	if p.word.tp == Tnum {
		p.word.val = sig + p.word.val // push signal to the beginning of the number
		n := newNode(p.word)
		p.next()
		return n
	}
	return p.Fail()
}
