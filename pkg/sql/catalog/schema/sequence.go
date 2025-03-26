package schema

import (
	"encoding/json"
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
)

func (s *Schema) LoadSequences(sequencesBytes [][]byte) error {
	for _, sequenceBytes := range sequencesBytes {
		var sequence *desc.Sequence
		err := json.Unmarshal(sequenceBytes, &sequence)
		if err != nil {
			return err
		}

		err = s.AddSequence(sequence)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) AddSequence(t *desc.Sequence) error {
	key := t.Name
	if _, ok := s.Sequences[key]; ok {
		return fmt.Errorf("table '%s' already exists", t.Name)
	}
	s.Sequences[key] = t
	return nil
}

func (s *Schema) GetSequence(key string) (*desc.Sequence, bool) {
	table, ok := s.Sequences[key]
	return table, ok
}

// DropSequence attempts to drop the table with the given key.
// If the table does not exist, it returns an error.
func (s *Schema) DropSequence(key string) error {
	_, ok := s.Sequences[key]
	if !ok {
		return fmt.Errorf("could not delete table '%s'", key)
	}
	delete(s.Sequences, key)
	return nil
}
