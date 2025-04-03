package db

import (
	"fmt"
	"os"
	"sync"

	"github.com/angles-n-daemons/popsql/pkg/db/executor"
	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/kv/store"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser"
	"github.com/angles-n-daemons/popsql/pkg/sql/plan"
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
	plan, err := plan.PlanQuery(e.Catalog.Schema, stmt)
	if err != nil {
		return err
	}

	exec := executor.New(e.Store, e.Catalog)
	rows, err := exec.Execute(plan)
	if err != nil {
		return err
	}
	fmt.Println(rows)
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
		if config.DebugPlanner {
			plan.Debug = true
		}
		db = newEngine(config.DebugStore)
	})
	return db
}
