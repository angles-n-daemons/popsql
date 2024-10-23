package kv

import "github.com/angles-n-daemons/popsql/pkg/kv/memtable"

var NewMemstore = memtable.NewMemstore

/*
Register is a savable object in the KV space.
*/
type Register interface {
	// primary key index
	ID() string
	Value() ([]byte, error)
	IndexIDs() ([]string, error)
}
