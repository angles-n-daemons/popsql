package ast

import (
	"fmt"
	"strings"

	"github.com/angles-n-daemons/popsql/pkg/debug/tree"
)

func Print(stmt Stmt) {
	t, err := VisitStmt(stmt, &stmtTreeifier{verbose: true, querifier: &ExprQuerifier{}})
	if err != nil {
		fmt.Println("Error in treeifier:", err)
		return
	}
	fmt.Println(tree.Visualize(t))
}

// - stmt treeifier
// - trickier with ast?
type stmtTreeifier struct {
	verbose   bool
	querifier *ExprQuerifier
}

func (t *stmtTreeifier) VisitCreateTableStmt(stmt *CreateTable) (*tree.Node, error) {
	content := []string{"CREATE TABLE: " + stmt.Name.Name.Lexeme}
	if t.verbose {
		for _, col := range stmt.Columns {
			content = append(content, fmt.Sprintf(" - %s %s", col.Name.Name.Lexeme, col.DataType.Lexeme))
		}
	}
	return tree.NewNode(content), nil
}

func (t *stmtTreeifier) VisitSelectStmt(stmt *Select) (*tree.Node, error) {
	content := []string{"SELECT: "}
	if stmt.From != nil {
		content[0] += stmt.From.Name.Lexeme
	}
	if t.verbose {
		terms := " terms: ["
		termsArr := []string{}
		for _, term := range stmt.Terms {
			ts, err := VisitExpr(term, t.querifier)
			if err != nil {
				return nil, err
			}
			termsArr = append(termsArr, ts)
		}
		terms += strings.Join(termsArr, ", ") + "]"

		if stmt.Where != nil {
			fs, err := VisitExpr(stmt.Where, t.querifier)
			if err != nil {
				return nil, err
			}
			filters := " filters: " + fs
			content = append(content, terms, filters)
		}
	}
	return tree.NewNode(content), nil
}

func (t *stmtTreeifier) VisitInsertStmt(stmt *Insert) (*tree.Node, error) {
	content := []string{"INSERT: " + stmt.Table.Name.Lexeme}
	if t.verbose {
		cols := " cols: ["
		colsArr := []string{}
		for _, col := range stmt.Columns {
			ts, err := VisitExpr(col, t.querifier)
			if err != nil {
				return nil, err
			}
			colsArr = append(colsArr, ts)
		}
		cols += strings.Join(colsArr, ", ") + "]"

		content = append(content, cols)
		content = append(content, " values: [")
		for _, tup := range stmt.Values {
			tupArr := []string{}
			for _, exp := range tup {
				ts, err := VisitExpr(exp, t.querifier)
				if err != nil {
					return nil, err
				}
				tupArr = append(tupArr, ts)
			}
			content = append(content, "  "+strings.Join(tupArr, ", "))

		}
		content = append(content, " ]")
	}
	return tree.NewNode(content), nil
}
