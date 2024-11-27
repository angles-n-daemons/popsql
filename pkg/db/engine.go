package db

import (
	"sync"

	"github.com/angles-n-daemons/popsql/pkg/kv/data"
	"github.com/angles-n-daemons/popsql/pkg/kv/memtable"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
	"github.com/angles-n-daemons/popsql/pkg/sys/catalog"
)

type Engine struct {
	Store   data.Store
	Catalog *catalog.Manager
}

func (e *Engine) Query(query string, parameters []any) error {
	stmt, err := parser.Parse(query)
	if err != nil {
		return err
	}
	switch stmt.(type) {
	case ast.Create:
		return e.Create(stmt)
	case ast.Insert:
		return nil
	}
	return nil
}

func newEngine() *Engine {
	store := memtable.NewMemstore()
	return &Engine{store, nil}
}

var db *Engine
var once sync.Once

func GetEngine() *Engine {
	once.Do(func() {
		db = newEngine()
	})
	return db
}
