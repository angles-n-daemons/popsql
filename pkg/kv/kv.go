package kv

import "github.com/angles-n-daemons/popsql/pkg/kv/keys"

/*
Register is a savable object in the KV space.
*/
type Register interface {
	// primary key index
	Key() *keys.Key
	Value() ([]byte, error)
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
