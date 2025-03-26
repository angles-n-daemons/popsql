package desc

import (
	"encoding/json"
	"strconv"
)

type Sequence struct {
	ID   uint64
	Name string
	V    uint64
}

func NewSequence(name string) *Sequence {
	return &Sequence{Name: name}
}

// Equal returns true if the two sequences are equal.
func (s *Sequence) Equal(o *Sequence) bool {
	return s.ID == o.ID && s.Name == o.Name && s.V == o.V
}

// Next increments the sequence value and returns the new value.
func (s *Sequence) Next() uint64 {
	s.V++
	return s.V
}

// Utility functions for the schema table
func (t *Sequence) Key() string {
	return strconv.FormatUint(t.ID, 10)
}

func (t *Sequence) Value() ([]byte, error) {
	return json.Marshal(t)
}
