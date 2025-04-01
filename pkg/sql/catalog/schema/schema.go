package schema

import "github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"

// The Schema struct is the in-memory representation of the database's schema.
// It exists so that the database can quickly access the definition of tables
// value of sequences among other things.
type Schema struct {
	Tables    *Collection[*desc.Table]
	Sequences *Collection[*desc.Sequence]
}

func NewSchema() *Schema {
	return &Schema{
		Tables:    NewCollection[*desc.Table](),
		Sequences: NewCollection[*desc.Sequence](),
	}
}

func SchemaFromCollections(
	tables *Collection[*desc.Table], sequences *Collection[*desc.Sequence],
) *Schema {
	return &Schema{
		Tables:    tables,
		Sequences: sequences,
	}
}

func Add[V desc.Any[V]](s *Schema, v V) error {
	return getCollection[V](s).Add(v)
}

func Get[V desc.Any[V]](s *Schema, id uint64) V {
	return getCollection[V](s).Get(id)
}

func GetByName[V desc.Any[V]](s *Schema, name string) V {
	return getCollection[V](s).GetByName(name)
}

func Remove[V desc.Any[V]](s *Schema, id uint64) error {
	return getCollection[V](s).Remove(id)
}

func getCollection[V desc.Any[V]](s *Schema) *Collection[V] {
	var zero V
	switch any(zero).(type) {
	case *desc.Table:
		return any(s.Tables).(*Collection[V])
	case *desc.Sequence:
		return any(s.Sequences).(*Collection[V])
	default:
		// this seems a little dangerous
		return nil
	}
}

// Empty returns whether the schema has no tables or sequences.
func Empty(s *Schema) bool {
	return s.Tables.Empty() && s.Sequences.Empty()
}

// Equal is a simple comparator which says whether two schema references
// are logically equivalent. It does this by checking whether the references'
// internal maps are equivalent in size, and whether the values for each of
// their keys are equivalent.
func (s *Schema) Equal(o *Schema) bool {
	if o == nil {
		return false
	}
	if !s.Tables.Equal(o.Tables) {
		return false
	}
	if !s.Sequences.Equal(o.Sequences) {
		return false
	}
	return true
}
