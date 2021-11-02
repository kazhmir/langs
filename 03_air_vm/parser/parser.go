package parser

import (
	. "air/parser/state"
	. "air/ast"
	"air/ast/T"
	"os"
	"fmt"
)

// Program = { Block }.
func Program(st *State) *Node {
	p := Repeat(st, Block)
	return &Node{
		NType: T.Program,
		Leafs: p,
	}
}

// Block = id "{" {Instr} "}" Branching ".".
func Block(st *State) *Node {
	if st.Word.NType != T.Id {
		return nil
	}
	name := Consume(st)
	Expect(st, T.LEFTBRACKET)
	instrs := &Node{
		NType: T.Code,
		Leafs: Repeat(st, Instr),
	}
	Expect(st, T.RIGHTBRACKET)
	jump := Branching(st)
	Expect(st, T.DOT)
	return &Node{
		NType: T.Block,
		Leafs: []*Node{name, instrs, jump},
	}
}

// Instr = id ":" Type Operand {"," Operand} ["->" id].
func Instr(st *State) *Node {
	if st.Word.NType != T.Id {
		return nil
	}
	name := Expect(st, T.Id)
	Expect(st, T.COLON)
	tp := Type(st)
	operands := &Node{Leafs: RepeatList(st, Operand, isComma)}
	var out *Node
	if st.Word.NType == T.ARROW {
		Consume(st)
		out = Expect(st, T.Id)
	}
	return &Node{
		NType: T.Instr,
		Leafs: []*Node{name, tp, operands, out},
	}
}
 
// Branching = If | GoTo | Ret.
func Branching(st *State) *Node {
	switch st.Word.NType {
		case T.GOTO:
			return GoTo(st)
		case T.IF:
			return If(st)
		case T.RET:
			return Ret(st)
	}
	ErrBadSymbol(st, "goto", "if", "ret")
	return nil
}
 
// GoTo = "goto" id.
func GoTo(st *State) *Node {
	kw := Expect(st, T.GOTO)
	id := Expect(st, T.Id)
	kw.Leafs = []*Node{id}
	return kw
}

// Ret = "ret" [Operand].
func Ret(st *State) *Node {
	kw := Expect(st, T.RET)
	if st.Word.NType != T.DOT {
		op := Operand(st)
		kw.Leafs = []*Node{op}
		return kw
	}
	return kw
}

// If = "if" ":" Type Operand "?" SwOption {"," SwOption}.
func If(st *State) *Node {
	kw := Expect(st, T.IF)
	Expect(st, T.COLON)
	tp := Type(st)
	op := Operand(st)
	Expect(st, T.INTERROGATION)
	branches := &Node{
		NType: T.Branches,
		Leafs: RepeatList(st, SwOption, isComma),
	}
	kw.Leafs = []*Node{tp, op, branches}
	return kw
}

// SwOption = Literal "->" id.
func SwOption(st *State) *Node {
	lit := Expect(st, T.Int, T.String, T.Bool)
	arrow := Expect(st, T.ARROW)
	id := Expect(st, T.Id)
	arrow.Leafs = []*Node{lit, id}
	return arrow
}
 
// Operand = id | Literal.
func Operand(st *State) *Node {
	switch st.Word.NType {
		case T.Id, T.String, T.Int, T.Bool:
			return Consume(st)
	}
	return nil
}

func Type(st *State) *Node {
	tp := Expect(st, T.Id)
	switch tp.Text {
		case "int":
			tp.NType = T.Int
		case "string":
			tp.NType = T.String
		case "bool":
			tp.NType = T.Bool
		default:
			fmt.Printf("%s Unknown type: %s\n", Place(st), tp.Text)
			os.Exit(1)
	}
	return tp
}

// Literal = num | string | bool.
// ---------------- isTk

func isComma(n *Node) bool {
	return n.NType == T.COMMA
}

