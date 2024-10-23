package keys

import "strings"

type Key struct {
	Table string
	ID    string
}

func NewKey(table string) *Key {
	return newKey(table, "")
}

func newKey(table, id string) *Key {
	return &Key{
		Table: table,
		ID:    id,
	}
}

func (k *Key) WithTable(table string) *Key {
	return newKey(table, k.ID)
}

func (k *Key) WithID(id string) *Key {
	return newKey(k.Table, id)
}

func (k *Key) WithIDAddition(id string) *Key {
	return newKey(k.Table, k.ID+id)
}

func (k *Key) String() string {
	key := k.Table
	if k.ID != "" {
		key += "/" + k.ID
	}
	return key
}

func NextString(s string) string {
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
