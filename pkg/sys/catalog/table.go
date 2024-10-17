package catalog

import (
	"fmt"
	"strings"
)

func NewTable(namespace *Namespace, name string, pkey []string) (*Table, error) {
	return nil, nil
}

type Table struct {
	Namespace  *Namespace
	Name       string
	Columns    []*Column
	PrimaryKey []string
}

func (t *Table) AddColumn(column *Column) error {
	return nil
}

func (t *Table) GetColumn(name string) (*Column, error) {
	for _, column := range t.Columns {
		if column.Name == name {
			return column, nil
		}
	}
	return nil, fmt.Errorf("Unable to find column %s on table %s", name, t.Name)
}

func (t *Table) Prefix() string {
	return fmt.Sprintf("%s/table/%s", t.Namespace.Name, t.Name)
}

func (t *Table) PrefixEnd() string {
	prefix := t.Prefix()
	return nextString(prefix)
}

func (t *Table) ToRegister() *TableRegister {
	return &TableRegister{
		t.Namespace.Name,
		t.Name,
		t.PrimaryKey,
	}
}

type TableRegister struct {
	Namespace  string
	Name       string
	PrimaryKey []string
}

func (t *TableRegister) Key() (string, error) {
	return t.Namespace + "-" + t.Name, nil
}

func (t *TableRegister) Value() ([]byte, error) {
	return nil, nil
}

func nextString(s string) string {
	i := len(s) - 1
	for i >= 0 && s[i] == 'z' {
		i--
	}

	if i == -1 {
		return s + "a"
	}

	j := 0
	return strings.Map(func(r rune) rune {
		if j == i {
			r += 1
		}
		j++
		return r
	}, s)
}
