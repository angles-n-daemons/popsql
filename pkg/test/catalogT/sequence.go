package catalogT

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
)

// Global counter to ensure each test-created table has a unique ID and (if needed) a unique name.
var sequenceIDCounter uint64

func SequenceID() uint64 {
	sequenceIDCounter++
	return sequenceIDCounter
}

func Sequence() *desc.Sequence {
	return NewSequence(nil)
}

func SequenceWithName(name string) *desc.Sequence {
	return NewSequence(&desc.Sequence{Name: name})
}

// Testing utility, which takes any portional part of a table and fills it out.
func NewSequence(s *desc.Sequence) *desc.Sequence {
	if s == nil {
		s = &desc.Sequence{}
	}

	if s.ID == 0 {
		s.ID = SequenceID()
	}

	if s.Name == "" {
		s.Name = fmt.Sprintf("sequence_%d", s.ID)
	}

	return s
}

func CopySequence(s *desc.Sequence) *desc.Sequence {
	return &desc.Sequence{
		ID:   s.ID,
		Name: s.Name,
		V:    s.V,
	}
}
