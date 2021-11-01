package T

type InstrType int

const (
	Add InstrType = iota
	Sub
	Div
	Mult
)

type Type int

const (
	Num Type = iota
	String
	Bool
)

type FlowType int

const (
	GOTO FlowType = iota
	IF
)
