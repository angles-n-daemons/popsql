package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/db/sql/parser"
	"github.com/angles-n-daemons/popsql/pkg/db/sql/parser/ast"
)

func TestParserBasic(t *testing.T) {
	stmt, err := parser.Parse(`SELECT * FROM users`)
	if err != nil {
		t.Fatal(err)
	}
	ast.Print(stmt)
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
		`SELECT jim`,
		`SELECT jim`,
		`SELECT jane, jane, jeffrey`,
		`SELECT * FROM users;`,
		`SELECT x, y FROM thing WHERE x==8`,
		`INSERT INTO a VALUES (1, 2)`,
		`INSERT INTO a (x, y) VALUES (1, 2)`,
		`INSERT INTO a (x, y) VALUES (1, 2), (3, 4)`,
		`CREATE TABLE derp()`,
		`CREATE TABLE derp (cal number)`,
		`CREATE TABLE derp(i string, cal number)`,
		`INSERT INTO JERP VALUES (1, 2)`,
		`INSERT INTO JERP (i, cal) VALUES (1, 2)`,
		//`UPDATE a SET x = 4`,
		//`UPDATE a SET x = 4, y = 5`,
		//`UPDATE a SET x = 4, y = 5 WHERE z = 10`,
		//`DELETE FROM a`,
		//`DELETE FROM a WHERE x=3`,
		// `DROP TABLE derp`,
		//`SELECT * FROM (SELECT * FROM b)`
	} {
		t.Run(`Parse Valid: `+query, func(t *testing.T) {
			fmt.Println("test in : ", query)
			stmt, err := parser.Parse(query)
			if err != nil {
				t.Fatal(`unexpected error for query: `, query, err)
			}
			s, err := ast.GenQuery(stmt)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println("test out: ", strings.Replace(s, "\n", "", -1))
			fmt.Println()
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
			_, err := parser.Parse(query)
			if err == nil {
				t.Fatalf(`expected query "%s" to fail parse, but didn"t`, query)
			}
		})
	}
}

// tests to run
// - ends with string
// - ends with number
// - ends with name
