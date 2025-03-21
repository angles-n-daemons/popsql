package schema

import (
	"encoding/json"
	"fmt"
)

const SchemaTableName = "__schema__"

type Schema struct {
	Sequences map[string]*Sequence
	Tables    map[string]*Table
}

func NewSchema() *Schema {
	schema := &Schema{
		Tables:    map[string]*Table{},
		Sequences: map[string]*Sequence{},
	}
	return schema
}

func (s *Schema) LoadTables(tablesBytes [][]byte) error {
	for _, tableBytes := range tablesBytes {
		var table *Table
		err := json.Unmarshal(tableBytes, &table)
		if err != nil {
			return err
		}
		s.AddTable(table)
	}

	return nil
}

func (s *Schema) AddTable(t *Table) error {
	key := t.Name
	if _, ok := s.Tables[key]; ok {
		return fmt.Errorf("table '%s' already exists", t.Name)
	}
	s.Tables[key] = t
	return nil
}

func (s *Schema) GetTable(key string) (*Table, error) {
	table, ok := s.Tables[key]
	if !ok {
		return nil, fmt.Errorf("could not find table '%s'", key)
	}
	return table, nil
}

// DropTable attempts to drop the table with the given key.
// If the table does not exist, it returns an error.
func (s *Schema) DropTable(key string) error {
	_, ok := s.Tables[key]
	if !ok {
		return fmt.Errorf("could not delete table '%s'", key)
	}
	delete(s.Tables, key)
	return nil
}

func (s *Schema) AddSequence(t *Sequence) error {
	key := t.Name
	if _, ok := s.Sequences[key]; ok {
		return fmt.Errorf("table '%s' already exists", t.Name)
	}
	s.Sequences[key] = t
	return nil
}

func (s *Schema) GetSequence(key string) (*Sequence, error) {
	table, ok := s.Sequences[key]
	if !ok {
		return nil, fmt.Errorf("could not find table '%s'", key)
	}
	return table, nil
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

// Equal is a simple comparator which says whether two schema references
// are logically equivalent. It does this by checking whether the references'
// internal maps are equivalent in size, and whether the values for each of
// their keys are equivalent.
func (s *Schema) Equal(o *Schema) bool {
	// if only one is nil, they cannot be equivalent
	if o == nil {
		return false
	}
	// if their internal maps are different sizes, they are not equivalent
	if len(s.Tables) != len(o.Tables) {
		return false
	}
	for key, table := range s.Tables {
		if !table.Equal(o.Tables[key]) {
			return false
		}
	}

	// if their internal maps are different sizes, they are not equivalent
	if len(s.Sequences) != len(o.Sequences) {
		return false
	}
	for key, sequence := range s.Sequences {
		if !sequence.Equal(o.Sequences[key]) {
			return false
		}
	}
	return true
}
