package ast

import (
	"fmt"
	"strings"
)

func printIndent(s string, indent int) {
	fmt.Println(strings.Repeat("\t", indent) + s)
}

type StmtPrinter struct {
	depth int
}

func PrintStmt(stmt Stmt) {
	VisitStmt(stmt, &StmtPrinter{})
}

func (p *StmtPrinter) print(stmt Stmt) {
	VisitStmt(stmt, p)
}

func (p *StmtPrinter) VisitSelectStmt(stmt *Select) (*any, error) {
	printIndent("- SELECT", p.depth)
	p.depth++
	printIndent("   terms:", p.depth-1)
	for _, term := range stmt.Terms {
		(&ExprPrinter{p.depth}).print(term)
	}
	printIndent("   from:", p.depth-1)
	if stmt.From != nil {
		(&ExprPrinter{p.depth}).print(stmt.From)
	}
	printIndent("  where:", p.depth-1)
	if stmt.From != nil {
		(&ExprPrinter{p.depth}).print(stmt.Where)
	}
	p.depth--
	return nil, nil
}

type ExprPrinter struct {
	depth int
}

func PrintExpr(expr Expr) {
	VisitExpr(expr, &ExprPrinter{})
}

func (p *ExprPrinter) print(expr Expr) {
	VisitExpr(expr, p)
}

func (p *ExprPrinter) VisitBinaryExpr(expr *Binary) (*any, error) {
	printIndent("- Binary", p.depth)
	p.depth++
	printIndent(fmt.Sprintf("     op: %s", expr.Operator.Type), p.depth-1)
	printIndent("   left:", p.depth-1)
	p.print(expr.Left)
	printIndent("  right:", p.depth-1)
	p.print(expr.Right)
	p.depth--
	return nil, nil
}

func (p *ExprPrinter) VisitLiteralExpr(expr *Literal) (*any, error) {
	printIndent(fmt.Sprintf("- Literal {%v}", expr.Value.Literal), p.depth)
	return nil, nil
}

func (p *ExprPrinter) VisitUnaryExpr(expr *Unary) (*any, error) {
	printIndent("- Unary", p.depth)
	p.depth++
	printIndent(fmt.Sprintf("   operator: %s", expr.Operator.Type), p.depth-1)
	printIndent("      right:", p.depth-1)
	p.print(expr.Right)
	p.depth--
	return nil, nil
}

func (p *ExprPrinter) VisitAssignmentExpr(expr *Assignment) (*any, error) {
	p.depth++
	p.print(expr.Value)
	p.depth--
	return nil, nil
}

func (p *ExprPrinter) VisitReferenceExpr(expr *Reference) (*any, error) {
	names := []string{}
	for _, token := range expr.Names {
		names = append(names, token.Lexeme)
	}
	printIndent(fmt.Sprintf("- Reference: %s", strings.Join(names, ".")), p.depth)
	return nil, nil
}
