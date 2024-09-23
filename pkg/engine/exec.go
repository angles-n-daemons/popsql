package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/angles-n-daemons/popsql/pkg/data"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
	"github.com/angles-n-daemons/popsql/pkg/sys"
)

func (db *Engine) VisitSelectStmt(stmt *ast.Select) (*any, error) {
	// read rows from table
	// ignore where clauses and what not for now
	// rearrange rows into specified column order
	// omit hidden columns
	table, err := db.lookupTable(stmt.From)
	if err != nil {
		return nil, err
	}
	rows, err := db.Store.Scan(table.KeyPrefix(), table.KeyPrefixEnd())
	if err != nil {
		return nil, err
	}
	records := []*sys.Record{}
	for _, row := range rows {
		record, err := sys.NewRecordFromBytes(table, stmt.Terms, row)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
		results := []string{}
		for _, val := range record.Data {
			results = append(results, fmt.Sprintf("%v", val))
		}
		fmt.Println(strings.Join(results, "|"))
	}
	return nil, nil
}

func (db *Engine) VisitInsertStmt(stmt *ast.Insert) (*any, error) {
	table, err := db.lookupTable(stmt.Table)
	if err != nil {
		return nil, err
	}
	for _, tuple := range stmt.Values {
		record, err := sys.NewRecordFromExpression(table, stmt.Columns, tuple)
		if err != nil {
			return nil, err
		}
		err = record.MaybeAddAutogenKey()
		key := fmt.Sprintf("%s/%s", table.KeyPrefix(), record.Key())
		b, err := record.ToBytes()
		if err != nil {
			return nil, err
		}
		err = db.Store.Put(key, b, false)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (db *Engine) VisitCreateStmt(stmt *ast.Create) (*any, error) {
	columns := []sys.Column{}
	name := stmt.Name.Lexeme
	for _, column := range stmt.Columns {
		dataType, err := sys.GetDataType(column.DataType)
		if err != nil {
			return nil, err
		}
		column := sys.Column{
			Space:    sys.USER,
			Table:    name,
			Name:     column.Name.Lexeme,
			DataType: dataType,
		}
		columns = append(columns, column)
	}
	db.CreateTable(sys.Table{Space: sys.USER, Name: name, Columns: columns})
	return nil, nil
}

func (db *Engine) CreateTable(table sys.Table) error {
	// add primary key if none set
	if len(table.PrimaryKey) == 0 {
		table.Columns = append(table.Columns, sys.Column{
			Space:    table.Space,
			Table:    table.Name,
			Name:     sys.PRIMARY_KEY_NAME,
			DataType: sys.STRING,
		})
		table.PrimaryKey = []string{sys.PRIMARY_KEY_NAME}
	}

	tablePrefix := table.KeyPrefix()
	tableKey := Key(db.Schema.System.Tables, &table)
	bytes, err := json.Marshal(table)
	if err != nil {
		return err
	}

	err = db.Store.Put(tableKey, bytes, false)
	if err != nil {
		return err
	}

	for _, column := range table.Columns {
		key := Key(db.Schema.System.Columns, &column)
		bytes, err = json.Marshal(table)
		if err != nil {
			return err
		}
		err = db.Store.Put(key, bytes, false)
		if err != nil {
			return err
		}
	}
	// verify it doesn't already exist
	db.Schema.Tables[tablePrefix] = table

	return nil
}

func (db *Engine) lookupTable(ref *ast.Reference) (*sys.Table, error) {
	if len(ref.Names) > 2 {
		return nil, fmt.Errorf("woah there cowboy, can't look up a table with %d names", len(ref.Names))
	}
	if len(ref.Names) == 0 {
		return nil, fmt.Errorf("woah there cowboy, can't look up a table with %d names", len(ref.Names))
	}
	names := []string{}
	if len(ref.Names) == 1 {
		names = append(names, sys.USER)
	}
	for _, nameExpr := range ref.Names {
		names = append(names, nameExpr.Lexeme)
	}
	key := fmt.Sprintf("%s/table/%s", names[0], names[1])
	table, ok := db.Schema.Tables[key]
	if !ok {
		return nil, fmt.Errorf("unable to find table %s", strings.Join(names, "."))
	}
	return &table, nil
}

func Key(table sys.Table, register sys.Register) string {
	return fmt.Sprintf("%s/%s", table.KeyPrefix(), register.Key())
}

func LoadSchema(store data.Store) (*sys.Schema, error) {
	schema := sys.NewSchema()
	tablesBytes, err := store.Scan(
		schema.System.Tables.KeyPrefix(),
		schema.System.Tables.KeyPrefixEnd(),
	)
	// ignore system tables
	if err != nil {
		return nil, err
	}

	for _, tableBytes := range tablesBytes {
		var table sys.Table
		err = json.Unmarshal(tableBytes, &table)
		if err != nil {
			return nil, err
		}
		// prevent override of table
		schema.Tables[table.KeyPrefix()] = table
	}
	columnsBytes, err := store.Scan(
		schema.System.Columns.KeyPrefix(),
		schema.System.Columns.KeyPrefixEnd(),
	)
	for _, columnBytes := range columnsBytes {
		var column sys.Column
		err = json.Unmarshal(columnBytes, &column)
		if err != nil {
			return nil, err
		}
		table, ok := schema.Tables[column.TableKeyPrefix()]
		//consistency in error
		if !ok {
			return nil, fmt.Errorf("unknown table for column: %s", column.TableKeyPrefix())
		}
		table.Columns = append(table.Columns, column)
	}
	return schema, nil
}

func ValidateReferences() error {
	return nil
}
