package catalog

import (
	"fmt"
	"strings"
)

type Table struct {
	Space      string
	Name       string
	Columns    []Column `json:"-"`
	PrimaryKey []string
}

func (t *Table) GetColumn(name string) (*Column, error) {
	for _, column := range t.Columns {
		if column.Name == name {
			return &column, nil
		}
	}
	return nil, fmt.Errorf("Unable to find column %s on table %s", name, t.Name)
}

func (t *Table) Prefix() string {
	return fmt.Sprintf("%s/table/%s", t.Space, t.Name)
}

func (t *Table) PrefixEnd() string {
	prefix := t.Prefix()
	return nextString(prefix)
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
