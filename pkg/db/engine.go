package db

import (
	"os"
	"sync"

	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/kv/store"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
)

type Engine struct {
	Store   kv.Store
	Catalog *catalog.Manager
}

func (e *Engine) Query(query string, parameters []any) error {
	stmt, err := parser.Parse(query)
	if err != nil {
		return err
	}
	switch v := stmt.(type) {
	case *ast.CreateTable:
		return e.CreateTable(v)
	case *ast.Insert:
		return nil
	}
	return nil
}

func newEngine(debugStore bool) *Engine {
	var st kv.Store = store.NewMemStore()
	if debugStore {
		st = store.NewDebugStore(st)
	}
	manager, err := catalog.NewManager(st)
	if err != nil {
		panic(err)
	}
	return &Engine{st, manager}
}

var db *Engine
var once sync.Once

func GetEngine() *Engine {
	once.Do(func() {
		config := NewConfig(os.Getenv)
		if config.DebugParser {
			parser.Debug = true
		}
		db = newEngine(config.DebugStore)
	})
	return db
}
