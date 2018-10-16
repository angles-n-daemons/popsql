package vm

import "github.com/angles-n-daemons/popsql/internal/core/vm/op"

type Instruction struct {
	op.OPCode
	P1 int32
	P2 int32
	P3 int32
	P4 interface{}
	P5 interface{}
}
