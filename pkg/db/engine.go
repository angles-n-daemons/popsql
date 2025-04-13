package db

import (
	"os"
	"sync"

	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/kv/store"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/sql/execution"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/sql/plan"
)

type Engine struct {
	Store   kv.Store
	Catalog *catalog.Manager
}

func (e *Engine) Query(query string, parameters []any) (*execution.Result, error) {
	stmt, err := parser.Parse(query)
	if err != nil {
		return nil, err
	}

	plan, err := plan.PlanQuery(e.Catalog.Schema, stmt)
	if err != nil {
		return nil, err
	}

	return execution.Run(e.Store, e.Catalog, plan)
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
		if config.DebugScanner {
			scanner.Debug = true
		}
		if config.DebugParser {
			parser.Debug = true
		}
		if config.DebugPlanner {
			plan.Debug = true
		}
		if config.DebugStore {
			desc.DebugTables = true
		}
		db = newEngine(config.DebugStore)
	})
	return db
}
