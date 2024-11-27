package kv

/*
Register is a savable object in the KV space.
*/
type Register interface {
	// primary key index
	ID() string
	Value() ([]byte, error)
	IndexIDs() ([]string, error)
}

type Store interface {
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error
	GetRange(start, end string) (Cursor, error)
}

type Cursor interface {
	ReadAll() ([][]byte, error)
	Read(num int) ([][]byte, error)
	IsAtEnd() bool
}
