package parser_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
)

func TestParserBasic(t *testing.T) {
	stmt, err := parser.Parse(`SELECT 'hi' == 'no' != 'where'`)
	if err != nil {
		t.Fatal(err)
	}
	ast.PrintStmt(stmt)
}

func TestValidPrograms(t *testing.T) {
	for _, query := range []string{
		"SELECT 1",
		"SELECT 1.23",
		"SELECT 'hi there'",
		"SELECT jim",
		"SELECT jim.jane",
		"SELECT jane.goodall, 12.3, 'jeremy lin'",
		//	"SELECT * FROM jim.jane",
		//	"select * from jim.jane",
		//	"select 'cal', col, 123. from jim.jane",
	} {
		stmt, err := parser.Parse(query)
		if err != nil {
			t.Fatal(err)
		}
		ast.PrintStmt(stmt)
	}
}

var invalidPrograms = []string{
	"SELECT 1;",
	"SELECT 1;",
}

// tests to run
// - ends with string
// - ends with number
// - ends with name
