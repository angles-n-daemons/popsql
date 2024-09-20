package engine

import (
	"os"

	"github.com/angles-n-daemons/popsql/pkg/data"
	"github.com/angles-n-daemons/popsql/pkg/data/memtable"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
	"github.com/angles-n-daemons/popsql/pkg/sys"
)

type Options struct {
	Filename string
}

type Engine struct {
	File   *os.File
	Store  data.Store
	Schema *sys.Schema
	New    bool
}

func NewEngine(opts Options) *Engine {
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

func (db *Engine) Query(query string) error {
	stmt, err := parser.Parse(query)
	if err != nil {
		return err
	}
	_, err = ast.VisitStmt(stmt, db)
	if err != nil {
		return err
	}
	return nil
}

func (db *Engine) CreateSystemTables() {
	db.CreateTable(db.Schema.System.Tables)
	db.CreateTable(db.Schema.System.Columns)
}
