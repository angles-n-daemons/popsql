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
Binary     = Expr Left, scanner.Token Operator, Expr Right
Literal    = scanner.Token Value
Unary      = scanner.Token Operator, Expr Right
Assignment = scanner.Token Name, Expr Value
Reference  = []*scanner.Token Names
`

var stmtAST = `
Select     = []Expr Terms, *Reference From, Expr Where
`

var walkFuncSignature = `
type walkFunc func(Expr) error
`

var exprInterface = `
type Expr interface { }
`

var stmtInterface = `
type Stmt interface { }
`

var errNilStr = `
%[1]sif err != nil {
%[1]s	return err
%[1]s}
`

var disclaimer = `
// this file was generated by tools/gen_ast.go
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

func parseGrammar(grammarStr string) []treeType {
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

	newline("package ast")

	newline(disclaimer)

	newline("import (")
	newline("\t\"fmt\"")
	newline("\t\"github.com/angles-n-daemons/popsql/pkg/sql/parser/scanner\"")
	newline(")")

	newline(walkFuncSignature)

	newline(exprInterface)
	newline(formatGrammar(exprAST, "Expr"))

	newline(formatGrammar(stmtAST, "Stmt"))
	newline(stmtInterface)

	return ast
}

func formatGrammar(grammarStr string, grammarType string) string {
	grammar := parseGrammar(grammarStr)
	output := ""
	newline := func(s string) {
		output += s + "\n"
	}

	newline(formatVisitor(grammar, grammarType))
	newline("")
	newline(formatVisit(grammar, grammarType))

	for _, ttype := range grammar {
		newline("")
		newline(formatStruct(ttype))
		newline("")
	}
	return output
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

func formatVisit(grammar []treeType, grammarType string) string {
	visitStr := fmt.Sprintf("func Visit%[1]s[T any](expr %[1]s, visitor %[1]sVisitor[T]) (*T, error) {\n", grammarType)
	newline := func(s string) {
		visitStr += s + "\n"
	}
	newline(fmt.Sprintf("\tswitch typed%s := expr.(type) {", grammarType))
	for _, ttype := range grammar {
		newline(fmt.Sprintf("\tcase *%s:", ttype.name))
		newline(fmt.Sprintf("\t\treturn visitor.Visit%s%[2]s(typed%[2]s)", ttype.name, grammarType))
	}
	newline("\tdefault:")
	newline(fmt.Sprintf("\t\treturn nil, fmt.Errorf(\"unable to visit type %%T\", typed%s)", grammarType))

	newline("\t}")
	newline("}")
	return visitStr
}

func formatVisitor(grammar []treeType, grammarType string) string {
	visitorStr := fmt.Sprintf("type %sVisitor[T any] interface {\n", grammarType)
	for _, ttype := range grammar {
		visitorStr += fmt.Sprintf("\tVisit%[1]s%[2]s(*%[1]s) (*T, error)\n", ttype.name, grammarType)
	}
	visitorStr += "}"
	return visitorStr
}
