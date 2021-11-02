package T

type NodeType int

const (
	Undefined NodeType = iota
	Block
	Instr
	Operand
	Variable
	Program
	Code
	Branches

	Int
	String
	Bool
	Id

	GOTO
	IF
	RET

	LEFTBRACKET
	RIGHTBRACKET
	ARROW
	COMMA
	INTERROGATION
	COLON
	DOT

	EOF
)

var tostring = map[NodeType]string{
	Undefined:     "Undefined ",
	Block:         "Block",
	Instr:         "Instr",
	Operand:       "Operand",
	Variable:      "Variable",
	Int:           "int",
	Id:            "id",
	String:        "String",
	Bool:          "Bool",
	GOTO:          "goto",
	IF:            "if",
	LEFTBRACKET:   "LEFTBRACKET",
	RIGHTBRACKET:  "RIGHTBRACKET",
	ARROW:         "ARROW",
	COMMA:         "COMMA",
	INTERROGATION: "INTERROGATION",
	COLON:         "COLON",
	DOT:           "DOT",
	EOF:           "EOF",
}

func FmtType(a ...NodeType) string {
	out := tostring[a[0]]
	for _, v := range a[1:] {
		out += ", " + tostring[v]
	}
	return out
}
