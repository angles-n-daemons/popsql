package ast

import (
	"fmt"
	"strings"
)

// queryify.go is a utility package for turning AST trees back
// into SQL query strings. It's useful for debugging a parse
// and testing whether the parsed out structure is similar to
// the original incoming query.
type StmtQuerifier struct {
	depth int
}

var stmtQuerifier = &StmtQuerifier{}
var exprQuerifier = &ExprQuerifier{}

func withIndent(depth int) string {
	return strings.Repeat("\t", depth)
}

func GenQuery(stmt Stmt) (string, error) {
	return VisitStmt(stmt, stmtQuerifier)
}

func (p *StmtQuerifier) toQuery(stmt Stmt) {
	VisitStmt(stmt, p)
}

func (p *StmtQuerifier) VisitCreateTableStmt(stmt *CreateTable) (string, error) {
	var sb strings.Builder
	w := sb.WriteString
	w(withIndent(p.depth) + "CREATE TABLE ")
	w(stmt.Name.Name.Lexeme + " (\n")
	colStrings := []string{}
	for _, col := range stmt.Columns {
		colStr, err := exprQuerifier.toQuery(col)
		if err != nil {
			return "", err
		}
		colStrings = append(colStrings, withIndent(p.depth+1)+colStr)
	}
	w(strings.Join(colStrings, ",\n"))
	w(withIndent(p.depth) + ")")
	s := sb.String()
	return s, nil
}

func (p *StmtQuerifier) VisitSelectStmt(stmt *Select) (string, error) {
	var sb strings.Builder
	w := sb.WriteString
	w(withIndent(p.depth) + "SELECT ")

	termStrings := []string{}
	for _, term := range stmt.Terms {
		termStr, err := exprQuerifier.toQuery(term)
		if err != nil {
			return "", err
		}
		termStrings = append(termStrings, termStr)
	}
	w(strings.Join(termStrings, ", ") + "\n")

	if stmt.From != nil {
		w(withIndent(p.depth) + "  FROM ")
		w(stmt.From.Name.Lexeme + "\n")
	}

	if stmt.Where != nil {
		w(withIndent(p.depth) + " FROM ")
		whereStr, err := exprQuerifier.toQuery(stmt.Where)
		if err != nil {
			return "", err
		}
		w(whereStr + "\n")
	}
	s := sb.String()
	return s, nil
}

func (p *StmtQuerifier) VisitInsertStmt(stmt *Insert) (string, error) {
	var sb strings.Builder
	w := sb.WriteString
	w(withIndent(p.depth) + "INSERT INTO ")
	w(stmt.Table.Name.Lexeme + "\n")
	columns := []string{}
	for _, column := range stmt.Columns {
		columns = append(columns, column.Name.Lexeme)
	}
	w("(" + strings.Join(columns, ", ") + ")")
	w(" VALUES ")

	tuples := []string{}
	for _, tuple := range stmt.Values {
		values := []string{}
		for _, val := range tuple {
			valStr, err := exprQuerifier.toQuery(val)
			if err != nil {
				return "", err
			}
			values = append(values, valStr)
		}
		tuples = append(tuples, "("+strings.Join(values, ", ")+")")
	}
	w(strings.Join(tuples, ", "))
	s := sb.String()
	return s, nil
}

type ExprQuerifier struct {
	depth int
}

func (p *ExprQuerifier) toQuery(expr Expr) (string, error) {
	return VisitExpr[string](expr, p)
}

func (p *ExprQuerifier) VisitIdentifierExpr(expr *Identifier) (string, error) {
	return expr.Name.Lexeme, nil
}

func (p *ExprQuerifier) VisitBinaryExpr(expr *Binary) (string, error) {
	left, err := p.toQuery(expr.Left)
	if err != nil {
		return "", err
	}

	right, err := p.toQuery(expr.Right)
	if err != nil {
		return "", err
	}

	s := fmt.Sprintf("%s %s %s", left, expr.Operator.Lexeme, right)
	return s, nil
}

func (p *ExprQuerifier) VisitLiteralExpr(expr *Literal) (string, error) {
	var s string
	switch v := expr.Value.Literal.(type) {
	case float64, bool:
		s = fmt.Sprintf(`%v`, v)
	case string:
		s = fmt.Sprintf(`"%s"`, v)
	default:
		return "", fmt.Errorf("unknown type %T", expr.Value.Literal)
	}
	return s, nil
}

func (p *ExprQuerifier) VisitUnaryExpr(expr *Unary) (string, error) {
	valueStr, err := p.toQuery(expr.Right)
	if err != nil {
		return "", err
	}
	s := expr.Operator.Lexeme + valueStr
	return s, nil
}

func (p *ExprQuerifier) VisitColumnSpecExpr(expr *ColumnSpec) (string, error) {
	s := expr.Name.Name.Lexeme + " " + expr.DataType.Lexeme
	return s, nil
}
