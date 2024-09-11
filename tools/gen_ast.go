package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(GenAST())
}

// The syntax tree grammar doesn't map to the actual grammar
// Eg where there's binary -> unary -> primary -> identifier
// each of those would be nodes, but in the ast binary just
// dips directly to the identifier, rather it just says it has
// some

var exprAST = `
Binary   = Expr left, Token operator, Expr right
`
var x = `
Grouping = Expr expr
Literal  = any value
Unary    = Token operator, Expr right
List     = []Expr exprs
`

var walkFuncSignature = `
type walkFunc[T any] func(Expr[T]) error
`

var exprInterface = `
type Expr[T any] interface {
	Accept(visitor Visitor[T]) error
	Walk(walkFunc) error
}
`

var errNilStr = `
%[1]sif err != nil {
%[1]s	return err
%[1]s}
`

type treeType struct {
	name   string
	fields []field
}

type field struct {
	name    string
	ftype   string
	isarray bool
}

func errNilReturn(indent int) string {
	pad := strings.Repeat("\t", indent)
	return fmt.Sprintf(errNilStr, pad)
}

func parseExprGrammar(grammarStr string) []treeType {
	types := []treeType{}
	typesStr := strings.Split(grammarStr, "\n")
	for _, ttype := range typesStr {
		if ttype == "" {
			continue
		}
		typeParts := strings.Split(ttype, "=")
		name := strings.TrimSpace(typeParts[0])
		props := strings.Split(strings.TrimSpace(typeParts[1]), ",")
		fields := []field{}
		for _, prop := range props {
			propParts := strings.Split(strings.TrimSpace(prop), " ")
			name := strings.TrimSpace(propParts[1])
			ftype := strings.TrimSpace(propParts[0])
			isarray := ftype[:2] == "[]"
			fields = append(fields, field{name, ftype, isarray})
		}
		types = append(types, treeType{name, fields})
	}
	return types

}

func GenAST() string {
	ast := ""
	newline := func(s string) {
		ast += s + "\n"
	}

	newline("package grammar")
	newline(walkFuncSignature)
	newline(exprInterface)

	grammar := parseExprGrammar(exprAST)

	newline(formatVisitor(grammar))

	for _, ttype := range grammar {
		newline("")
		newline(formatStruct(ttype))
		newline("")
		newline(formatWalk(ttype))
		newline("")
		newline(formatVisit(ttype))
		newline("")
	}
	return ast
}

func formatStruct(ttype treeType) string {
	structStr := fmt.Sprintf("type %s struct {\n", ttype.name)
	newline := func(s string) {
		structStr += s + "\n"
	}
	for _, field := range ttype.fields {
		newline(fmt.Sprintf("\t%s %s", field.name, field.ftype))
	}
	newline("}")
	return structStr
}

func formatWalk(ttype treeType) string {
	walkStr := fmt.Sprintf("func (e %s) Walk(f walkFunc) error {\n", ttype.name)
	newline := func(s string) {
		walkStr += s + "\n"
	}
	newline("\tvar err error")
	for _, field := range ttype.fields {
		if field.isarray {
			newline(fmt.Sprintf("\tfor i := 0; i < len(e.%s); i++ {", field.name))
			newline(fmt.Sprintf("\t\terr = e.%s[i].Walk(f)", field.name))
			newline(errNilReturn(2))
			newline("\t}")
		} else {
			newline(fmt.Sprintf("\terr = e.%s.Walk(f)", field.name))
			newline(errNilReturn(1))
		}
	}
	newline("\treturn nil")
	newline("}")
	return walkStr
}

func formatVisit(ttype treeType) string {
	visitStr := fmt.Sprintf("func (e %s) Visit(v ExprVisitor) error {\n", ttype.name)
	visitStr += fmt.Sprintf("\treturn v.Visit%[1]sExpr(e)\n", ttype.name)
	visitStr += "}"
	return visitStr
}

func formatVisitor(grammar []treeType) string {
	visitorStr := "type ExprVisitor interface {\n"
	for _, ttype := range grammar {
		visitorStr += fmt.Sprintf("\tVisit%[1]sExpr(%[1]s) error\n", ttype.name)
	}
	visitorStr += "}"
	return visitorStr
}
