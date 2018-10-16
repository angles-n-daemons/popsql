package vm

import "github.com/angles-n-daemons/popsql/internal/core/vm/op"

func NewVM() *VM {
	// Return a VM with 100 registers
	return &VM{
		registers: make([]interface{}, 100),
		counter:   0,
	}
}

type VM struct {
	registers []interface{}
	counter   int32
	ExitCode  int32
}

func (vm *VM) Run(program []*Instruction) ([][]interface{}, error) {
	results := make([][]interface{}, 0)

	for {
		// fetch/decode the instruction
		ins := program[vm.counter]

		// execute the instruction
		output, err := vm.exec(ins)
		if err != nil {
			return nil, err
		}

		// if output given add to result set
		if output != nil {
			results = append(results, output)
		}

		// if next instruction is -1 exit
		if vm.counter == -1 {
			break
		}
	}

	return results, nil
}

func (vm *VM) exec(ins *Instruction) ([]interface{}, error) {
	var output []interface{}
	nextInstruction := vm.counter + 1
	var err error

	switch ins.OPCode {
	case op.Halt:
		// set the exit code
		vm.ExitCode = ins.P1
		// set nextInstruction to -1 so program exits
		nextInstruction = -1
		// TODO: use P2 to rollback transaction

	case op.Init:
		// Go to instruction P2
		nextInstruction = ins.P2
		// Increment P1 to show program has entered
		ins.P1++
		// TODO: use address P3 for SQLITE_CORRUPT error

	case op.Integer:
		// Write integer P1 into register P2
		vm.registers[ins.P2] = ins.P1

	case op.ResultRow:
		// return the registers between P1 and P2 + 1
		output = vm.registers[ins.P1 : ins.P2+1]

	}

	vm.counter = nextInstruction
	return output, err
}
