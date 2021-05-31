package re

import (
	"log"
)

type nodeType int

const (
	and nodeType = iota
	or
	star
	set
	emptyStr
)

var nodeTypePrint = map[nodeType]string{
	and:      "and",
	or:       "or",
	star:     "star",
	set:      "set",
	emptyStr: "empty string",
}

type node struct {
	set      *Set
	tp       nodeType
	children []*node
}

/*
this prints a sideways view of the Abstract Syntax Tree
*/
func (n node) String() string {
	return n.beautify(0)
}

func (n node) beautify(d int) string {
	output := "{" + nodeTypePrint[n.tp] + "}\n"
	if n.set != nil {
		output = "{" + n.set.String() + ":" + nodeTypePrint[n.tp] + "}\n"
	}
	for _, child := range n.children {
		output += indent(d) + child.beautify(d+1)
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

type parser struct {
	inp  []token
	word token
	path string
	i    int
}

func (p *parser) next() {
	if p.i < len(p.inp) {
		p.word = p.inp[p.i]
		p.i++
	}
}

func (p *parser) unread() {
	p.i--
}

func (p *parser) run(s []token) *node {
	if len(s) == 0 {
		return &node{set: NewSet("", false)}
	}
	p.i = 0
	p.inp = s
	p.path = "run:" + p.word.String() + "\n"
	p.next()
	root := p.expr()
	return root
}

func (p *parser) expr() *node {
	p.path += "expr:" + p.word.String() + "\n"
	n := p.str()
	if n == nil {
		return n
	}
	if p.word.val == '|' && p.word.tp == ope {
		leafs := []*node{n}
		for p.word.val == '|' && p.word.tp == ope {
			p.next()
			next := p.str()
			if next != nil {
				leafs = append(leafs, next)
			}
		}
		return &node{
			tp:       or,
			children: leafs,
		}
	}
	return n
}

func (p *parser) str() *node {
	p.path += "str:" + p.word.String() + "\n"
	a := p.rep()
	if b := p.rep(); b != nil {
		leafs := []*node{a, b}
		for n := p.rep(); n != nil; n = p.rep() {
			leafs = append(leafs, n)
		}
		return &node{
			tp:       and,
			children: leafs,
		}
	}
	return a
}

func (p *parser) rep() *node {
	p.path += "rep:" + p.word.String() + "\n"
	n := p.term()
	if n == nil {
		return nil
	}
	if p.word.val == '*' && p.word.tp == ope {
		p.next()
		return &node{
			tp:       star,
			children: []*node{n},
		}
	}
	return n
}

func (p *parser) term() *node {
	p.path += "term:" + p.word.String() + "\n"
	if p.word.val == '(' && p.word.tp == ope {
		p.next()
		n := p.expr()
		p.next() // discards (
		return n
	}
	if p.word.val == '[' && p.word.tp == ope {
		p.next()
		return p.set()
	}
	if p.word.tp == char { // can be optimized
		r := p.word.val
		p.next()
		return &node{
			set: NewSet(string(r), false),
			tp:  set,
		}
	}
	if p.word.tp == empty {
		p.next()
		return &node{
			tp: emptyStr,
		}
	}
	// operator
	return nil
}

func (p *parser) set() *node {
	p.path += "set:" + p.word.String() + "\n"
	out := &node{set: NewSet("", false), tp: set}
	if p.word.val == '^' && p.word.tp == ope {
		out.set.Negated = true
		p.next()
	}

	items := map[rune]struct{}{}
	// the only operator expected here is ']'
	for ; p.word.tp != ope; p.next() {
		items = union(items, p.item())
	}
	if p.word.val == ']' {
		p.next() // discards ']'
	} else {
		log.Fatalf("unexpected operator %c in set", p.word.val)
	}

	out.set.Items = items
	return out
}

func (p *parser) item() map[rune]struct{} {
	p.path += "item:" + p.word.String() + "\n"
	first := p.word.val
	p.next()
	if p.word.val == '-' && p.word.tp == ope {
		p.next()
		if p.word.val != ']' && p.word.tp != ope {
			return runeRange(first, p.word.val)
		}
		log.Fatal("Range operator requires two operands")
		return nil
	}
	p.unread()
	return map[rune]struct{}{
		first: struct{}{},
	}
}

func runeRange(a, b rune) map[rune]struct{} {
	if a > b {
		hold := a
		a = b
		b = hold
	}
	output := make(map[rune]struct{}, b-a)
	for c := a; c <= b; c++ {
		output[c] = struct{}{}
	}
	return output
}
