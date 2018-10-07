package vm

type Instruction struct {
	OPCode
	P1 interface{}
	P2 interface{}
	P3 interface{}
	P4 interface{}
	P5 interface{}
}

type OPCode int

const (
	Abortable OPCode = iota + 1
	Add
	AddImm
	Affinity
)
