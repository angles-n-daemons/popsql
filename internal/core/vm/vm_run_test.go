package vm

import (
	"reflect"
	"testing"

	"github.com/angles-n-daemons/popsql/internal/core/vm/op"
)

func TestSelectThree(t *testing.T) {
	vm := NewVM()
	program := []*Instruction{
		// start at instruction 1
		{
			OPCode: op.Init,
			P2:     1,
		},
		// set register 1 to be 3
		{
			OPCode: op.Integer,
			P1:     3,
			P2:     1,
		},
		// add result row from registers 1 to 1
		{
			OPCode: op.ResultRow,
			P1:     1,
			P2:     1,
		},
		// halt the program
		{
			OPCode: op.Halt,
			P1:     0,
		},
	}

	result, err := vm.Run(program)
	if err != nil {
		t.Fatal(err)
	}

	expected := [][]interface{}{
		{int32(3)},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf(
			"Run expected %v got %v",
			expected,
			result,
		)
	}
}
