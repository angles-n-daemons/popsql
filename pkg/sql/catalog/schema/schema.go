package schema

import "github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"

const SchemaTableName = "__schema__"

type Schema struct {
	Sequences map[string]*desc.Sequence
	Tables    map[string]*desc.Table
}

func NewSchema() *Schema {
	schema := &Schema{
		Tables:    map[string]*desc.Table{},
		Sequences: map[string]*desc.Sequence{},
	}
	return schema
}

// Empty returns whether the schema has no tables or sequences.
func (s *Schema) Empty() bool {
	return len(s.Tables) == 0 && len(s.Sequences) == 0
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
