package kv

import "unicode/utf8"

// The END_ID is used to denote the end of a table span.
const END_ID = '\u0000'

// MIN_RUNE skips the control characters
const MIN_RUNE = '\u0020'

// Key is the reference object for a key in the database's keyspace.
// It contains a table, the set of records the key belongs to and an ID,
// the identifier for the individual record.
type Key struct {
	Table string
	ID    string
}

// NewKey creates a new key from a table.
func NewKey(table string) *Key {
	return newKey(table, "")
}

// newKey creates a new key from a table and an id.
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
	id := k.ID
	if isEnd(k.ID) {
		id = "<END>"
	}
	key += "/" + id
	return key
}

func isEnd(id string) bool {
	return utf8.RuneCountInString(id) == 1 && []rune(id)[0] == END_ID
}

func (k *Key) Next() *Key {
	if k.ID == "" {
		return newKey(k.Table, string(END_ID))
	}
	return newKey(k.Table, NextString(k.ID))
}

func NextString(s string) string {
	// Can't be a next after the end.
	if isEnd(s) {
		return s
	}

	n := len(s)
	if n == 0 {
		return string(MIN_RUNE)
	}

	runes := []rune(s)
	if runes[len(runes)-1] == utf8.MaxRune || len(s) == 0 {
		runes = append(runes, MIN_RUNE)
	} else {
		runes[len(runes)-1]++
	}
	return string(runes)
}
