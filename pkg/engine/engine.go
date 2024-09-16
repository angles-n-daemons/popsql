package engine

import (
	"github.com/angles-n-daemons/popsql/pkg/sql/parser"
	"github.com/angles-n-daemons/popsql/pkg/sql/parser/ast"
)

type Engine struct {
	DataDir string
	// theoretically a lot of configuration
}

func (e *Engine) Query(query string) error {
	stmt, err := parser.Parse(query)
	if err != nil {
		return err
	}
	ast.PrintStmt(stmt)
	return nil
}
