package parser_test

import (
	"testing"

	"github.com/angles-n-daemons/popsql/pkg/sql/parser"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
)

func TestParserBasic(t *testing.T) {
	expr, err := parser.Parse(`'hi' == 'no' != 'where'`)
	if err != nil {
		t.Fatal(err)
	}
	(&ast.Printer{}).Print(expr)
}
