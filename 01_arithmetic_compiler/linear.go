package main

import (
	"fmt"
)

type Operator int

const (
	MOV Operator = iota
	OUT

	ADD
	MUL
	DIV
	SUB
)

var OpToStr = map[Operator]string{
	MOV: "MOV",
	OUT: "OUT",

	ADD: "ADD",
	MUL: "MUL",
	DIV: "DIV",
	SUB: "SUB",
}

var SymbToOp = map[string]Operator{
	"+": ADD,
	"*": MUL,
	"/": DIV,
	"-": SUB,
}

type OpType int

const (
	tREGI OpType = iota
	tADDR
	tNUMB
	tARGU
)

type Operand struct {
	Data string // either a register number, a number or address
	Type OpType
}

func (op *Operand) String() string {
	if op.Type == tREGI {
		return "R" + op.Data
	}
	return op.Data
}

type Instr struct {
	Op      Operator
	a, b, c *Operand
}

func (i *Instr) String() string {
	if i.c != nil {
		return fmt.Sprintf("%s %v, %v -> %v\n", OpToStr[i.Op], i.a, i.b, i.c)
	}
	if i.b != nil {
		return fmt.Sprintf("%s %v -> %v\n", OpToStr[i.Op], i.a, i.b)
	}
	return fmt.Sprintf("%s %v\n", OpToStr[i.Op], i.a)
}

type Block []*Instr

func (b Block) String() string {
	out := "block {\n"
	for i := range b {
		out += b[i].String()
	}
	out += "}\n"
	return out
}

type CodeGen struct {
	Counter int
	Code    *Block
}

func (cg *CodeGen) Generate(n *node) *Block {
	cg.Code = &Block{}
	cg.Counter = 0
	out := cg.gen(n)
	cg.AddInstr(&Instr{
		Op: OUT,
		a:  out,
	})
	return cg.Code
}

func (cg *CodeGen) gen(n *node) *Operand {
	switch n.tp {
	case Tnum:
		return &Operand{
			Data: n.val,
			Type: tNUMB,
		}
	case Tcla:
		return &Operand{
			Data: n.val,
			Type: tARGU,
		}
	}
	if len(n.leafs) == 1 { // unary
		out := cg.gen(n.leafs[0])
		if n.val == "-" { // MUL num, -1 -> Rx
			return cg.GenInverse(out) // returns Rx
		}
		return out
	}
	a := cg.gen(n.leafs[0])
	b := cg.gen(n.leafs[1])
	return cg.GenOp(n.val, a, b) // <OP> a, b -> Rx
}

func (cg *CodeGen) NextRegister() *Operand {
	r := cg.Counter
	cg.Counter++
	return &Operand{
		Data: fmt.Sprint(r),
	}
}

func (cg *CodeGen) AddInstr(i *Instr) {
	*cg.Code = append(*cg.Code, i)
}

func (cg *CodeGen) GenOp(op string, a, b *Operand) *Operand {
	out := cg.NextRegister()
	cg.AddInstr(&Instr{
		Op: SymbToOp[op],
		a:  a,
		b:  b,
		c:  out,
	})
	return out
}

func (cg *CodeGen) GenInverse(a *Operand) *Operand {
	b := &Operand{
		Data: "-1",
		Type: tNUMB,
	}
	out := cg.NextRegister()
	cg.AddInstr(&Instr{
		Op: MUL,
		a:  a,
		b:  b,
		c:  out,
	})
	return out
}
