package catalog

import "github.com/angles-n-daemons/popsql/pkg/sys/schema"

// Bootstrap is the entry point for a fresh database. If nothing is in the store,
// it will automatically populate the two system tables, the meta table and the
// sequences table, and create their corresponding sequences.
// It needs to do this however using the direct functions on the schema as well
// as the store, since normal operation depends on their already existing.
func (m *Manager) Bootstrap() error {
	err := m.bootstrapSequence(InitMetaTableSequence)
	if err != nil {
		return err
	}

	err = m.bootstrapSequence(InitSequencesTableSequence)
	if err != nil {
		return err
	}

	err = m.bootstrapTable(InitMetaTable)
	if err != nil {
		return err
	}

	err = m.bootstrapTable(InitSequencesTable)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) bootstrapSequence(s *schema.Sequence) error {
	err := m.Schema.AddSequence(s)
	if err != nil {
		return err
	}
	err = m.storeSequence(InitSequencesTable, s)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) bootstrapTable(t *schema.Table) error {
	err := m.Schema.AddTable(t)
	if err != nil {
		return err
	}
	err = m.storeTable(InitMetaTable, t)
	if err != nil {
		return err
	}
	return nil
}
