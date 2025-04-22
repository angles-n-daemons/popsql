package kv

/*
Register is a savable object in the KV space.
*/
type Register interface {
	// primary key index
	Key() string
	Value() ([]byte, error)
}

// Store is the primary interface which will be exposed to the
// rest of the database system. It defines a simple key-value
// interface, which can be implemented in any chosen way; in
// memory for example, as an LSM, or even with a B+ tree.
type Store interface {
	Get(string) ([]byte, error)
	Put(string, []byte) error
	Scan(start, end string) (Cursor, error)
}

// Cursor holds the results for a scan, to be read out in batches
// so that a caller can control the volume of data they're processing.
type Cursor interface {
	ReadAll() ([][]byte, error)
	Read(num int) ([][]byte, error)
	Next() ([]byte, error)
	IsAtEnd() bool
}
