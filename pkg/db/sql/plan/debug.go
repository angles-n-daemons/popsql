package plan

import "fmt"

type PlanDebugger struct {
	verbose bool
	depth   int
}

func DebugPlan(plan Plan) (string, error) {
	debugger := &PlanDebugger{}
	return VisitPlan(plan, debugger)
}

func (p *PlanDebugger) VisitCreateTable(plan *CreateTable) (string, error) {
	output := "CreateTable: " + plan.Table.Name()
	if p.verbose {
		output += "\n"
	}
	return "CreateTable: " + plan.Table.Name(), nil
}

func (p *PlanDebugger) VisitInsert(plan *Insert) (string, error) {
	return "Insert: " + plan.Table.Name(), nil
}

func (p *PlanDebugger) VisitScan(plan *Scan) (string, error) {
	return "Scan: " + plan.Table.Name(), nil
}

func (p *PlanDebugger) VisitValues(plan *Values) (string, error) {
	return fmt.Sprintf("Values: %d rows", len(plan.Rows)), nil
}
