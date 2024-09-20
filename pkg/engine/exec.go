package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/angles-n-daemons/popsql/pkg/data"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
	"github.com/angles-n-daemons/popsql/pkg/sys"
)

func (db *Engine) PrintRows(rows []any) {

}

func (db *Engine) VisitSelectStmt(stmt *ast.Select) (*any, error) {
	// read rows from table
	// ignore where clauses and what not for now
	// rearrange rows into specified column order
	// omit hidden columns
	return nil, nil
}

func (db *Engine) VisitInsertStmt(stmt *ast.Insert) (*any, error) {
	// lookup table for validation
	// insert to rows
	//   order rows in order of columns
	//   validate data types
	//   add missing column data
	//   reorder to match default column ordering
	// write rows
	table, err := db.lookupTable(stmt.Table)
	if err != nil {
		return nil, err
	}
	tupleLen := len(stmt.Columns)
	columns := make([]*sys.Column, tupleLen)
	for i, col := range stmt.Columns {
		column, err := table.GetColumn(col)
		if err != nil {
			return nil, err
		}
		columns[i] = column
	}
	for i, tuple := range stmt.Values {
		fmt.Println(tuple)
		if len(tuple) != tupleLen {
			return nil, fmt.Errorf("%dth value set found incorrect size %d (expected %d)", i, len(tuple), tupleLen)
		}
		var record map[string]any
		for j, value := range tuple {
			// validate data type
			// how to handle references?
			// only handle simple data types
			switch value.(type) {
			case *ast.Literal:
				fmt.Println(value, j)
			default:
				return nil, fmt.Errorf("unexpected expression type %T reading VALUES", value)
			}
		}
		// validate record
		fmt.Println(record, tuple)
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
	// validate primary key
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
	for tk, t := range db.Schema.Tables {
		fmt.Println(tk, t)
		fmt.Printf("same? '%s' '%s' %d %d\n", tk, key, len(tk), len(key))
	}
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
