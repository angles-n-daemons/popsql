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

// features
// - insert statements
// - update statements
// - delete statements
// - aliasing? (later)
// - group by? (later)
// - capitalization?
// - insert with select
// - condition with single equal
func TestValidPrograms(t *testing.T) {
	for _, query := range []string{
		"SELECT 1",
		"SELECT 1.23",
		"SeleCT 1.23",
		"sELEct 1.23",
		"SELECT 'hi there'",
		"SELECT jim",
		"SELECT jim.jane",
		"SELECT jane.goodall, 12.3, 'jeremy lin'",
		"SELECT * FROM users;",
		"SELECT column.* FROM users.thing",
		"SELECT column.*, 12 FROM users.thing",
		"SELECT 5 + 4, 'ello' FROM thing WHERE x==8",
		"INSERT INTO a (x, y) VALUES (1, 2)",
		"INSERT INTO a (x, y) VALUES (1, 2), (3, 4)",
		"INSERT INTO a.b (c.d) VALUES (5)",
		"SELECT !false",
		"SELECT a.false",
		"SELECT (1)",
	} {
		stmt, err := parser.Parse(query)
		if err != nil {
			t.Fatal("unexpected error for query: ", query, err)
		}
		ast.PrintStmt(stmt)
	}
}

func TestInvalidPrograms(t *testing.T) {
	for _, query := range []string{
		"SELECT",
		"5 + 4",
		"SELECT FROM 5+4",
		"SELECT * FROM",
		"SELECT * FROM z WHERE",
		"SELECT * FROM z SELECT *",
		"SELECT * FROM '",
		"SELECT * FROM 5",
		"SELECT * WHERE SELECT *",
		"INSERT",
		"INSERT FROM",
		"INSERT INTO 5",
		"INSERT INTO a.b thing",
		"INSERT INTO a.b (5",
		"INSERT INTO a.b (c.d(",
		"INSERT INTO a.b (c.d) FROM",
		"INSERT INTO a.b (c.d) VALUES 5",
		"INSERT INTO a.b (c.d) VALUES (5, ",
		"INSERT INTO a.b (c.d) VALUES (5 a",
		"SELECT !",
		"SELECT (5 + 4",
	} {
		stmt, err := parser.Parse(query)
		if err == nil {
			t.Fatalf("expected query '%s' to fail parse, but didn't", query)
		}
		ast.PrintStmt(stmt)
	}
}

// tests to run
// - ends with string
// - ends with number
// - ends with name
