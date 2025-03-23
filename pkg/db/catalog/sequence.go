package catalog

import (
	"fmt"

	"github.com/angles-n-daemons/popsql/pkg/sys/schema"
)

func (m *Manager) storeSequence(sequenceTable *schema.Table, s *schema.Sequence) error {
	key := sequenceTable.Prefix().WithID(s.Key())
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
