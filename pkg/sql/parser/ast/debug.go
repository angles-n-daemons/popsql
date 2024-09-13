package ast

import (
	"fmt"
	"strings"
)

type Printer struct {
	depth int
}

func (p *Printer) Print(expr Expr) {
	Visit(expr, p)
}

func printIndent(s string, indent int) {
	fmt.Println(strings.Repeat("\t", indent) + s)
}

func (p *Printer) VisitBinaryExpr(expr *Binary) (*any, error) {
	printIndent("- Binary", p.depth)
	p.depth++
	printIndent(fmt.Sprintf("   operator: %s", expr.Operator.Type), p.depth-1)
	printIndent("       left:", p.depth-1)
	p.Print(expr.Left)
	printIndent("      right:", p.depth-1)
	p.Print(expr.Right)
	p.depth--
	return nil, nil
}

func (p *Printer) VisitLiteralExpr(expr *Literal) (*any, error) {
	printIndent(fmt.Sprintf("--- Literal %v", expr.Value), p.depth)
	return nil, nil
}

func (p *Printer) VisitUnaryExpr(expr *Unary) (*any, error) {
	p.depth++
	p.Print(expr.Right)
	p.depth--
	return nil, nil
}

func (p *Printer) VisitAssignmentExpr(expr *Assignment) (*any, error) {
	p.depth++
	p.Print(expr.Value)
	p.depth--
	return nil, nil
}

func (p *Printer) VisitExprListExpr(expr *ExprList) (*any, error) {
	p.depth++
	for _, e := range expr.Exprs {
		p.Print(e)
	}
	p.depth--
	return nil, nil
}
