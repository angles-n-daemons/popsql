package vm

import (
	"reflect"
	"testing"

	"github.com/angles-n-daemons/popsql/internal/core/vm/op"
)

func TestInstructionInit(t *testing.T) {
	vm := NewVM()
	ins := Instruction{
		OPCode: op.Init,
		P1:     0,
		P2:     3,
	}

	_, err := vm.exec(&ins)
	if err != nil {
		t.Fatal(err)
	}

	if ins.P1 != 1 {
		t.Fatalf(
			"Expected P1 on instruction to be %d but got %d",
			1,
			ins.P1,
		)
	}

	if vm.counter != 3 {
		t.Fatalf(
			"Expected VM program counter to be %d but got %d",
			3,
			vm.counter,
		)
	}
}

func TestInstructionInteger(t *testing.T) {
	vm := NewVM()
	ins := Instruction{
		OPCode: op.Integer,
		P1:     1,
		P2:     2,
	}

	_, err := vm.exec(&ins)
	if err != nil {
		t.Fatal(err)
	}

	if vm.registers[ins.P2] != ins.P1 {
		t.Fatalf(
			"Expected register %d to be %d but instead got %d",
			ins.P2,
			ins.P1,
			vm.registers[ins.P2],
		)
	}
}

func TestInstructionResultRow(t *testing.T) {
	vm := NewVM()
	vm.registers = []interface{}{1, 2, 3, 4, 5}
	ins := Instruction{
		OPCode: op.ResultRow,
		P1:     1,
		P2:     3,
	}

	output, err := vm.exec(&ins)
	if err != nil {
		t.Fatal(err)
	}

	expected := []interface{}{2, 3, 4}

	if !reflect.DeepEqual(expected, output) {
		t.Fatalf(
			"ResultRow expected %v got %v",
			expected,
			output,
		)
	}
}

func TestInstructionHalt(t *testing.T) {
	vm := NewVM()
	ins := Instruction{
		OPCode: op.Halt,
		P1:     7,
	}

	_, err := vm.exec(&ins)
	if err != nil {
		t.Fatal(err)
	}

	if vm.counter != -1 {
		t.Fatalf(
			"Expected VM program counter to be %d but got %d",
			-1,
			vm.counter,
		)
	}

	if vm.ExitCode != 7 {
		t.Fatalf(
			"Expected VM exit code to be %d but got %d",
			7,
			vm.ExitCode,
		)
	}
}

func TestInstructionString8(t *testing.T) {
	vm := NewVM()
	ins := Instruction{
		OPCode: op.String8,
		P2:     1,
		P4:     "test",
	}

	_, err := vm.exec(&ins)
	if err != nil {
		t.Fatal(err)
	}

	if vm.registers[ins.P2] != "test" {
		t.Fatalf(
			"Expected register 2 to be %s but got %v",
			"test",
			vm.registers[ins.P2],
		)
	}
}
