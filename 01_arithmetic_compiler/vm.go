package main

import (
	"fmt"
	"strconv"
)

type Machine struct {
	Regs map[string]int
}

func (m *Machine) Run(code *Block) {
	m.Regs = make(map[string]int, 10)
	var a, b int
	for _, ins := range *code {
		if ins.a != nil {
			a = m.GetOperand(ins.a)
		}
		if ins.b != nil {
			b = m.GetOperand(ins.b)
		}
		switch ins.Op {
		case ADD:
			m.Regs[ins.c.Data] = a + b
		case SUB:
			m.Regs[ins.c.Data] = a - b
		case MUL:
			m.Regs[ins.c.Data] = a * b
		case DIV:
			m.Regs[ins.c.Data] = a / b
		case OUT:
			fmt.Println(a)
		}
	}
}

func (m *Machine) GetOperand(op *Operand) int {
	switch op.Type {
	case tNUMB:
		n, _ := strconv.ParseInt(op.Data, 10, 64)
		return int(n)
	case tREGI:
		n, ok := m.Regs[op.Data]
		if ok {
			return n
		}
		panic("Using a register that doesn't exist")
	case tARGU:
		index := LangArgToIndex(op.Data)
		n, _ := strconv.ParseInt(lang_args[index], 10, 64)
		return int(n)
	}
	panic("Unimplemented")
}
