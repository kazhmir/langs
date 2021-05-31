package main

import (
	"fmt"
	"os"
)

type node struct {
	*lexeme
	leafs []*node
}

func newNode(l *lexeme) *node {
	return &node{
		lexeme: l,
	}
}

func (n *node) String() string {
	return n.ast(0)
}

func (n *node) AddLeaf(l *node) {
	n.leafs = append(n.leafs, l)
}

func (n *node) ast(i int) string {
	output := n.lexeme.String() + "\n"
	for l := range n.leafs {
		output += indent(i) + n.leafs[l].ast(i+1)
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
		parent := newNode(p.word) // "+" or "-" node
		parent.AddLeaf(last)
		p.next()
		parent.AddLeaf(p.Term())
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
	last := p.Unary()
	for p.word.val == "*" || p.word.val == "/" {
		parent := newNode(p.word) // "*" or "/" node
		parent.AddLeaf(last)
		p.next()
		parent.AddLeaf(p.Unary())
		last = parent
	}
	return last
}

func (p *Parser) Unary() *node {
	if p.word.val == "-" || p.word.val == "+" {
		parent := newNode(p.word)
		p.next()
		parent.AddLeaf(p.Factor())
		return parent
	}
	return p.Factor()
}

func (p *Parser) Factor() *node {
	if p.word.val == "(" {
		p.next()
		return p.Expr()
	}
	return p.Num()
}

func (p *Parser) Num() *node {
	if p.word.tp == Tnum || p.word.tp == Targ {
		n := newNode(p.word)
		p.next()
		return n
	}
	return p.Fail()
}
