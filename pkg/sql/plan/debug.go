package plan

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
	return "Insert into: " + plan.Table.Name(), nil
}
