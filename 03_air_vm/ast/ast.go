package ast

import (
	"air/ast/T"
	"fmt"
)

type Node struct {
	Text  string
	NType T.NodeType

	Line int
	Col  int

	Leafs []*Node
}

func AddLeaf(a *Node, b *Node) {
	a.Leafs = append(a.Leafs, b)
}

func (n *Node) String() string {
	return n.ast(0)
}

func (n *Node) ast(i int) string {
	output := fmt.Sprintf("{'%s', %s}",
		n.Text,
		T.FmtType(n.NType),
	)
	for _, kid := range n.Leafs {
		if kid == nil {
			output += indent(i) + "nil"
			continue
		}
		output += indent(i) + kid.ast(i+1)
	}
	return output
}

func indent(n int) string {
	output := "\n"
	for i := -1; i < n-1; i++ {
		output += "   "
	}
	output += "└─>"
	return output
}
