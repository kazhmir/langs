package main

import "fmt"

/*
Physical registers are just a number between 0 and 13, if 14
registers are given as available.

Next[ physical ] returns the distance to next use of value
Location[ virtual ] retuns the register storing this value
Value[ virtual ] retuns the value stored in the register
*/
type Resources struct {
	Available *Stack // stack of free registers
	Next      []int  // distance to next ocurrence of register

	Location map[string]int // virtual register -> physical register
	Value    map[int]string // physical register -> virtual register
}

/* Free frees the given physical register */
func (r *Resources) Free(pReg int) {
	vReg := r.Value[pReg]
	delete(r.Value, pReg)
	delete(r.Location, vReg)
	r.Next[pReg] = 1 << 32 // a large number so it has higher priority to be allocated
	r.Available.Push(pReg)
}

/*FurthestUse finds the virtual register that's further from
the current instruction, observe that we don't need the current position
since the largest index would be the furthest away anyway*/
func (r *Resources) FurthestUse() int {
	dist := 0
	reg := 0
	for i := range r.Next {
		if r.Next[i] > dist {
			dist = r.Next[i]
			reg = i
		}
	}
	return reg
}

type Allocator struct {
	in   *Block
	curr int
	out  string

	Address   map[string]int // virtual register -> stack offset
	rbpOffset int            // current offset from top of stack
}

func (alc *Allocator) Begin(res *Resources) string {
	alc.out = Header
	for _, ins := range *alc.in {
		alc.curr++
		if ins.c != nil { // 3 operands, ADD, MUL, SUB, DIV
			if ins.Op == DIV {
				alc.GenDiv(ins, res)
				continue
			}
			regC := alc.Alloc(ins.c.Data, res)
			regA := alc.GenCode(ins.a, "mov", regC, res)
			regB := alc.GenCode(ins.b, OpToASM[ins.Op], regC, res)

			alc.FreeIfNotNeeded(ins.a, regA, res)
			alc.FreeIfNotNeeded(ins.b, regB, res)
			continue
		}
		if ins.b != nil { // 2 operands, always MOV
			pReg := alc.Alloc(ins.b.Data, res)
			alc.out += fmt.Sprintf("\tmov\t%s, %s\n", x64Reg[pReg], ins.a.Data)
			continue
		}
		pReg := alc.Ensure(ins.a, res)
		alc.out += fmt.Sprintf("\tpush\t%s\n", x64Reg[pReg])
	}
	alc.out += Tail
	return alc.out
}

/* GenCode ensures the value is either in a register or is a literal
then generates the code with pRegOut as the target register.
It's implemented as a separate function to avoid repetition for each operand
*/
func (alc *Allocator) GenCode(op *Operand, ins string, pRegOut int, res *Resources) int {
	switch op.Type {
	case tNUMB:
		alc.out += fmt.Sprintf("\t%s\t%s, %s\n", ins, x64Reg[pRegOut], op.Data)
	case tARGU, tREGI:
		pReg := alc.Ensure(op, res)
		alc.out += fmt.Sprintf("\t%s\t%s, %s\n", ins, x64Reg[pRegOut], x64Reg[pReg])
		return pReg
	}
	return pRegOut
}

/* FreeIfNotNeeded check's if the value is needed, and if not, frees the physical register
Literals are never needed, so it just returns.
*/
func (alc *Allocator) FreeIfNotNeeded(op *Operand, pReg int, res *Resources) {
	if op.Type == tNUMB {
		return
	}
	if !alc.IsNeeded(op.Data, pReg, res) {
		res.Free(pReg)
	}
}

/* Ensure ensures that the operand is in a register,
if the operand is a literal, it returns a invalid register (-1)
if a literal needs to be in a register for a particular operation
other code must be used (See Allocator.GenDiv method)
*/
func (alc *Allocator) Ensure(op *Operand, res *Resources) int {
	if v, ok := res.Location[op.Data]; ok {
		return v
	}
	pReg := alc.Alloc(op.Data, res)
	if op.Type == tARGU {
		offset := 16 + 8*LangArgToIndex(op.Data)
		alc.out += fmt.Sprintf("\tmov\t%s, [rbp + %v]\n", x64Reg[pReg], offset)
		return pReg
	}
	if v, ok := alc.Address[op.Data]; ok {
		alc.out += fmt.Sprintf("\tmov\t%s, [rbp - %v]\n", x64Reg[pReg], v)
		return pReg
	}
	panic("Something went wrong")
}

/* IsNeeded performs a linear scan through the input instructions
to see if the value is used again. If it finds one use, it updates
Resources.Next with the index of the next use and returns true.
*/
func (alc *Allocator) IsNeeded(vReg string, pReg int, res *Resources) bool {
	for i, ins := range (*alc.in)[alc.curr:] {
		if ins.a.Data == vReg || (ins.b != nil && ins.b.Data == vReg) {
			res.Next[pReg] = alc.curr + i
			return true
		}
	}
	return false
}

/* Alloc allocates a physical register to hold the value represented
by the virtual register. If no physical register is available, it finds
the physical register that's needed further from the current instruction index
and pushes the value to the stack (See Allocator.GenStore).
*/
func (alc *Allocator) Alloc(vReg string, res *Resources) int {
	pReg := -1
	if res.Available.IsEmpty() {
		pReg = res.FurthestUse() // find biggest value in Next
		val := res.Value[pReg]
		alc.GenStore(val, pReg) // generate store for the value in the register
		res.Free(pReg)          // push on top of stack
	}
	pReg = res.Available.Pop() // pop top of stack
	res.Location[vReg] = pReg
	res.Value[pReg] = vReg
	res.Next[pReg] = -1 // to guarantee it's not used by the current operation
	return pReg
}

/* GenStore generates code to push the given value to the stack,
it also saves the address as an offset from the base pointer for
further references to this virtual register.
The base pointer is set at the initialization as the very top of the
stack (mov rbp, rsp), the address [rbp] then stores the number of
console arguments given to the program.
*/
func (alc *Allocator) GenStore(vReg string, pReg int) {
	alc.out += "\tpush\t" + x64Reg[pReg] + "\n"
	alc.rbpOffset += 8
	alc.Address[vReg] = alc.rbpOffset
}

/*MoveOrStore moves the given value out of the physical register holding it,
if there's an available register it just moves, otherwise it generates a store.*/
func (alc *Allocator) MoveOrStore(vReg string, pReg int, res *Resources) {
	if res.Available.IsEmpty() {
		alc.GenStore(vReg, pReg)
		return
	}
	newpReg := alc.Alloc(vReg, res)
	alc.out += fmt.Sprintf("\tmov\t%s, %s\n", x64Reg[newpReg], x64Reg[pReg])
}

/* GenDiv generates code for a division. Since the idiv instruction
uses both RAX and RDX to store the quotient and remainder respectively,
it needs special care.
Also, RDX must be zeroed before the instruction, otherwise you'll get a
Floating point exception. Note that we don't use the value on RDX, we're just cleaning it,
so we don't need to allocate and therefore we don't need to free at the end.
*/
func (alc *Allocator) GenDiv(ins *Instr, res *Resources) {
	const rax, rdx = 0, 3
	alc.EspecEnsure(ins.a, rax, res) // ensure A is in rax, regardless of it's type
	alc.EnsureFree(rdx, res)         // ensure rdx is free

	alc.out += fmt.Sprintf("\txor\t%s, %s\n", x64Reg[rdx], x64Reg[rdx]) // make sure rdx is zeroed
	pReg := alc.Ensure(ins.b, res)
	if pReg < 0 { // ins.b is a number
		reg := alc.Alloc(ins.b.Data, res) // idiv requires it to be in a register
		alc.out += fmt.Sprintf("\tmov\t%s, %s\n", x64Reg[reg], ins.b.Data)
		alc.out += fmt.Sprintf("\tidiv\t%s\n", x64Reg[reg])
		res.Free(reg) // literals are not needed
	} else {
		alc.out += fmt.Sprintf("\tidiv\t%s\n", x64Reg[pReg])
		alc.FreeIfNotNeeded(ins.b, pReg, res)
	}
	res.Location[ins.c.Data] = rax
	res.Value[rax] = ins.c.Data
}

/* EspecEnsure ensures that the value is stored in a specific register
this is used in the idiv instruction, since it always explicitly
uses the RAX and RDX registers as the result
*/
func (alc *Allocator) EspecEnsure(op *Operand, pReg int, res *Resources) {
	vRegOld := alc.EnsureFree(pReg, res)
	if op.Type == tNUMB {
		alc.out += fmt.Sprintf("\tmov\t%s, %v\t\n", x64Reg[pReg], op.Data)
		return
	}
	// here op.Data is a virtual register
	if vRegOld == op.Data {
		return
	}
	if pRegOld, ok := res.Location[op.Data]; ok {
		alc.out += fmt.Sprintf("\tmov\t%s, %s\n", x64Reg[pReg], x64Reg[pRegOld])
		return
	}
	if offset, ok := alc.Address[op.Data]; ok {
		alc.out += fmt.Sprintf("\tmov\t%s, [rbp - %v]\n", x64Reg[pReg], offset)
		return
	}
	fmt.Printf("This shouldn't execute %#v\n", op)
}

/* EnsureFree ensures that the physical register is free,
if there's a needed value in it, it moves the value somewhere else.
*/
func (alc *Allocator) EnsureFree(pReg int, res *Resources) string {
	vRegOld, ok := res.Value[pReg]
	if ok && alc.IsNeeded(vRegOld, pReg, res) { // if the register is occupied and the value is needed
		alc.MoveOrStore(vRegOld, pReg, res)
	}
	return vRegOld
}

type Stack struct {
	Data []int
	Top  int
}

func NewStack(ln int) *Stack {
	data := make([]int, ln)
	for i := 0; i < ln; i++ {
		data[i] = i
	}
	return &Stack{
		Data: data,
		Top:  ln - 1,
	}
}

func (s *Stack) Push(r int) {
	s.Top++
	s.Data[s.Top] = r
}

func (s *Stack) Pop() int {
	v := s.Data[s.Top]
	s.Top--
	return v
}

func (s *Stack) IsEmpty() bool {
	if s.Top < 0 {
		return true
	}
	return false
}

func (s *Stack) IsFull() bool {
	if s.Top >= len(s.Data)-1 {
		return true
	}
	return false
}
