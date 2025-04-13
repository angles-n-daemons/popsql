package desc

import (
	"encoding/json"
	"strconv"
)

type Sequence struct {
	SID   uint64 `json:"id"`
	SName string `json:"name"`
	V     uint64 `json:"value"`
}

func NewSequence(name string) *Sequence {
	return &Sequence{SName: name}
}

func NewSequenceFromArgs(id uint64, name string, val uint64) *Sequence {
	return &Sequence{SID: id, SName: name, V: val}
}

func (s *Sequence) WithID(id uint64) {
	s.SID = id
}

func (s *Sequence) ID() uint64 {
	return s.SID
}

func (s *Sequence) Name() string {
	return s.SName
}

// Equal returns true if the two sequences are equal.
func (s *Sequence) Equal(o *Sequence) bool {
	if o == nil {
		return false
	}
	return s.SID == o.SID && s.SName == o.SName && s.V == o.V
}

// Next increments the sequence value and returns the new value.
func (s *Sequence) Next() uint64 {
	s.V++
	return s.V
}

// Utility functions for the schema table
func (t *Sequence) Key() string {
	return strconv.FormatUint(t.SID, 10)
}

func (t *Sequence) Value() ([]byte, error) {
	return json.Marshal(t)
}
