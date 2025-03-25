package catalog

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sys/schema/desc"
)

func (c *Catalog) SequenceNext(s *desc.Sequence) (uint64, error) {
	next := s.Next()

	err := c.storeSequence(s)
	if err != nil {
		return 0, err
	}
	return next, nil
}

func (c *Catalog) addSequence(s *desc.Sequence) (*desc.Sequence, error) {
	id, err := c.SequenceNext(c.sequencesTableSequence)
	if err != nil {
		return nil, err
	}
	s.ID = id

	err = c.Schema.AddSequence(s)
	if err != nil {
		return nil, err
	}

	err = c.storeSequence(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (c *Catalog) storeSequence(s *desc.Sequence) error {
	key := c.sequencesTable.Prefix().WithID(s.Key())
	sequenceBytes, err := s.Value()
	if err != nil {
		return fmt.Errorf("failed encoding sequence while saving to store %w", err)
	}
	err = c.Store.Put(key.Encode(), sequenceBytes)
	if err != nil {
		return fmt.Errorf("could not put sequence definition in store %w", err)
	}
	return nil
}
