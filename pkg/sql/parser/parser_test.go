package parser_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
)

func TestParserBasic(t *testing.T) {
	stmt, err := parser.Parse(`SELECT "hi" == "no" != "where"`)
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
func TestParseValidPrograms(t *testing.T) {
	for _, query := range []string{
		`SELECT
		1`,
		`SELECT 1`,
		`SELECT 1;`,
		`SELECT 1.23`,
		`SeleCT 1.23`,
		`sELEct 1.23`,
		`SELECT "hi there"`,
		`SELECT jim`,
		`SELECT jim.jane`,
		`SELECT jane.goodall, 12.3, "jeremy lin"`,
		`SELECT * FROM users;`,
		`SELECT column.* FROM users.thing`,
		`SELECT column.*, 12 FROM users.thing`,
		`SELECT 5 + 4, "ello" FROM thing WHERE x==8`,
		`INSERT INTO a VALUES (1, 2)`,
		`INSERT INTO a (x, y) VALUES (1, 2)`,
		`INSERT INTO a (x, y) VALUES (1, 2), (3, 4)`,
		`INSERT INTO a.b (c.d) VALUES (5)`,
		`SELECT !false`,
		`SELECT a.false`,
		`SELECT (1)`,
		`CREATE TABLE derp()`,
		`CREATE TABLE derp (cal number)`,
		`CREATE TABLE derp(i string, cal number)`,
		//`UPDATE a SET x = 4`,
		//`UPDATE a SET x = 4, y = 5`,
		//`UPDATE a SET x = 4, y = 5 WHERE z = 10`,
		//`DELETE FROM a`,
		//`DELETE FROM a WHERE x=3`,
		// `DROP TABLE derp`,
	} {
		t.Run(`Parse Valid: `+query, func(t *testing.T) {
			stmt, err := parser.Parse(query)
			if err != nil {
				t.Fatal(`unexpected error for query: `, query, err)
			}
			ast.PrintStmt(stmt)
		})
	}
}

func TestParseInvalidPrograms(t *testing.T) {
	for _, query := range []string{
		`SELECT`,
		`5 + 4`,
		`SELECT FROM 5+4`,
		`SELECT * FROM`,
		`SELECT * FROM z WHERE`,
		`SELECT * FROM z SELECT *`,
		`SELECT * FROM "`,
		`SELECT * FROM 5`,
		`SELECT * WHERE SELECT *`,
		`INSERT`,
		`INSERT FROM`,
		`INSERT INTO 5`,
		`INSERT INTO a.b thing`,
		`INSERT INTO a.b (5`,
		`INSERT INTO a.b (c.d(`,
		`INSERT INTO a.b (c.d) FROM`,
		`INSERT INTO a.b (c.d) VALUES 5`,
		`INSERT INTO a.b (c.d) VALUES (5, `,
		`INSERT INTO a.b (c.d) VALUES (5 a)`,
		`INSERT INTO (c.d) a.b VALUES (5)`,
		`SELECT !`,
		`SELECT (5 + 4`,
		`CREATE TABLE x`,
	} {
		t.Run(`Parse Invalid: `+query, func(t *testing.T) {
			stmt, err := parser.Parse(query)
			if err == nil {
				t.Fatalf(`expected query "%s" to fail parse, but didn"t`, query)
			}
			ast.PrintStmt(stmt)
		})
	}
}

// tests to run
// - ends with string
// - ends with number
// - ends with name
