package keys

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

func (k *Key) Next() *Key {
	if k.ID == "" {
		return newKey(NextString(k.Table), k.ID)
	}
	return newKey(k.Table, NextString(k.ID))
}

func NextString(s string) string {
	n := len(s)
	if n == 0 {
		return "a"
	}

	runes := []rune(s)
	if runes[len(runes)-1] == 'z' || len(s) == 0 {
		runes = append(runes, 'a')
	} else {
		runes[len(runes)-1]++
	}
	return string(runes)
}
