package re

import (
	"log"
	"sort"
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
		return &node{set: &Set{Items: []rune{}}}
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
		log.Fatal("Wrong syntax for operator")
	}
	if p.word.val == '|' && p.word.tp == ope {
		leafs := []*node{n}
		for p.word.val == '|' && p.word.tp == ope {
			p.next()
			leafs = append(leafs, p.str())
		}
		if p.word.val == ')' && p.word.tp == ope {
			p.next()
		}
		return &node{
			tp:       or,
			children: leafs,
		}
	}
	if p.word.val == ')' && p.word.tp == ope {
		p.next()
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
		return p.expr()
	}
	if p.word.val == '[' && p.word.tp == ope {
		p.next()
		return p.set()
	}
	if p.word.tp == char { // can be optimized
		r := p.word.val
		p.next()
		return &node{
			set: &Set{
				Items: []rune{r},
			},
			tp: set,
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
	out := &node{set: &Set{Items: []rune{}}, tp: set}
	if p.word.val == '^' && p.word.tp == ope {
		out.set.Negated = true
		p.next()
	}

	s := runeSlice{}
	// the only operator expected here is ']'
	for ; p.word.tp != ope; p.next() {
		s = append(s, p.item()...)
	}
	if p.word.val == ']' {
		p.next() // discards ']'
	} else {
		log.Fatalf("unexpected operator %c in set", p.word.val)
	}

	ordered := rmDuplicates(s)
	sort.Sort(&ordered)
	out.set.Items = ordered
	return out
}

func (p *parser) item() runeSlice {
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
	return runeSlice{first}
}

func runeRange(a, b rune) []rune {
	if a > b {
		c := a
		a = b
		b = c
	}
	output := make([]rune, b-a)
	for i, c := 0, a; c <= b; i, c = i+1, c+1 {
		if i < len(output) {
			output[i] = c
			continue
		}
		output = append(output, c)
	}
	return output
}
