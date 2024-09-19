package engine

import (
	"encoding/json"
	"fmt"
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
	store := memtable.NewMemStore()
	isNew := true
	db := &Engine{Store: store, Schema: sys.NewSchema()}
	if isNew {
		db.CreateSystemTables()
	}
	db.LoadSchema()

	return db
}

func (db *Engine) Query(query string) error {
	stmt, err := parser.Parse(query)
	if err != nil {
		return err
	}
	ast.PrintStmt(stmt)
	ast.VisitStmt(stmt, db)
	db.LoadSchema()
	return nil
}

func (db *Engine) CreateSystemTables() {
	db.CreateTable(db.Schema.System.Tables)
	db.CreateTable(db.Schema.System.Columns)
}

func (db *Engine) CreateTable(table sys.Table) error {
	key := Key(db.Schema.System.Tables, &table)
	bytes, err := json.Marshal(table)
	if err != nil {
		return err
	}

	err = db.Store.Put(key, bytes, false)
	if err != nil {
		return err
	}

	for _, column := range table.Columns {
		key = Key(db.Schema.System.Columns, &column)
		bytes, err = json.Marshal(table)
		if err != nil {
			return err
		}
		err = db.Store.Put(key, bytes, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func Key(table sys.Table, register sys.Register) string {
	return fmt.Sprintf("%s/%s", table.KeyPrefix(), register.Key())
}

func (db *Engine) LoadSchema() (*sys.Schema, error) {
	tablesBytes, err := db.Store.Scan(
		db.Schema.System.Tables.KeyPrefix(),
		db.Schema.System.Tables.KeyPrefixEnd(),
	)
	// ignore system tables
	if err != nil {
		return nil, err
	}

	for _, tableBytes := range tablesBytes {
		fmt.Println(string(tableBytes))
		var table sys.Table
		err = json.Unmarshal(tableBytes, &table)
		if err != nil {
			return nil, err
		}

	}
	return nil, nil
}
