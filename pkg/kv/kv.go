package kv

type Store interface {
	Get(string) ([][]byte, error)
	// Gets range [start, end)
	GetRange(string, string) Cursor
	Put(string, []byte) error
	PutRange(string, string, []byte) error
}

/*
SELECT MANY
  - GetRange (store)

SELECT ONE
  - GetRange?

INSERT MANY
  - Put

INSERT ONE
  - Put

UPDATE MANY
  - GetRange and Put

DELETE MANY
  - DeleteRange
*/
type Cursor interface {
	Read(int) ([][]byte, error)
	ReadAll() ([][]byte, error)
	IsAtEnd() bool
}

/*
Register is a savable object in the KV space.
*/
type Register interface {
	// primary key index
	ID() string
	Value() ([]byte, error)
	IndexIDs() ([]string, error)
}
