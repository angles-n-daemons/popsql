package db

import (
	"sync"

	"github.com/angles-n-daemons/popsql/pkg/kv/data"
)

type Engine struct {
	Store  data.Store
	Schema *catalog.Schema
	New    bool
}

type Options struct {

}


func newEngine(opts Options) *Engine {
	// need some heavy debug flags
	// might be worth tagging the logging
	store := memtable.NewMemStore()
	isNew := true
	schema, err := LoadSchema(store)
	if err != nil {
		// do not panic here
		panic(err)
	}
	db := &Engine{Store: store, Schema: schema}
	if isNew {
		db.CreateSystemTables()
	}
	return db
}


var db *Engine
var once sync.Once

func GetEngine(opts Options) *Engine {
	once.Do(func() {
		db = newEngine(opts)
	})
	return db
}
