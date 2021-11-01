package air

import (
	"air/ir/T"
)

//type Function struct {
//	Name   string
//	Params []*types.TypeLayout
//	Start  *SimpleBlock
//}

type SimpleBlock struct {
	Name       string
	Operations []Instr
	Jump       Flow
}

type Instr struct {
	Tp        T.Type
	Operation T.InstrType
	Operands  []Operand
}

/*
GOTO has only Default
IF has multiple conditions and optinally a Default
*/
type Flow struct {
	Tp      T.FlowType
	Cond    []Condition
	Default *SimpleBlock
	Value   Operand
}

type Condition struct {
	Value  Operand
	Target *SimpleBlock
}

/*
Literals, variables, identifiers, types
depends on the operation (OpType)
*/
type Operand interface{}
