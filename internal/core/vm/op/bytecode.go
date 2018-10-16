package op

type OPCode int

const (
	Abortable OPCode = iota + 1
	Add
	AddImm
	Affinity
	Halt
	Init
	Integer
	ResultRow
)
