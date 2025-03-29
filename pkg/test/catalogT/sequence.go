package catalogT

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/kv"
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
	return NewSequence(&desc.Sequence{SName: name})
}

// Testing utility, which takes any portional part of a table and fills it out.
func NewSequence(s *desc.Sequence) *desc.Sequence {
	if s == nil {
		s = &desc.Sequence{}
	}

	if s.SID == 0 {
		s.SID = SequenceID()
	}

	if s.SName == "" {
		s.SName = fmt.Sprintf("sequence_%d", s.SID)
	}

	return s
}

func CopySequence(s *desc.Sequence) *desc.Sequence {
	return &desc.Sequence{
		SID:   s.SID,
		SName: s.SName,
		V:     s.V,
	}
}

func ReadSequence(t *testing.T, st kv.Store, key string) *desc.Sequence {
	sequenceBytes, err := st.Get(key)
	if err != nil {
		t.Fatal(t)
	}
	var s *desc.Sequence
	err = json.Unmarshal(sequenceBytes, &s)
	if err != nil {
		t.Fatal(t)
	}
	return s
}
