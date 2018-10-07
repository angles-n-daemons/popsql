package processor

import (
	"github.com/angles-n-daemons/popsql/internal/core/vm"
	"github.com/angles-n-daemons/popsql/internal/db"
)

func NewProcessor(c compiler) *Processor {
	return &Processor{c}
}

/*
Processor is the main entry point to the query parser and
virtual machine. It is responsible for processing queries.
*/
type Processor struct {
	compiler
}

func (p *Processor) Query(s string) (*db.Table, error) {
	_, err := p.compiler.Compile(s)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

type compiler interface {
	Compile(string) ([]vm.Instruction, error)
}
