package catalog

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
)

func (m *Manager) SequenceNext(s *desc.Sequence) (uint64, error) {
	// Get the next value in the sequence.
	next := s.Next()

	// Update the sequence in the store.
	err := m.StoreSequence(s)
	if err != nil {
		return 0, err
	}
	return next, nil
}

func (m *Manager) CreateSequence(s *desc.Sequence) (*desc.Sequence, error) {
	// create an id for the new sequence.
	id, err := m.SequenceNext(m.Sys.SequencesTableSequence)
	if err != nil {
		return nil, err
	}
	s.ID = id

	// save the sequence in the memory schema.
	err = m.Schema.AddSequence(s)
	if err != nil {
		return nil, err
	}

	err = m.StoreSequence(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (m *Manager) StoreSequence(s *desc.Sequence) error {
	key := m.Sys.SequencesTable.Prefix().WithID(s.Key())
	sequenceBytes, err := s.Value()
	if err != nil {
		return fmt.Errorf("failed encoding sequence while saving to store %w", err)
	}
	err = m.Store.Put(key.Encode(), sequenceBytes)
	if err != nil {
		return fmt.Errorf("could not put sequence definition in store %w", err)
	}
	return nil
}
