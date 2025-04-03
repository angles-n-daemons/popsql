package plan

type PlanDebugger struct{}

func DebugPlan(plan Plan) (string, error) {
	debugger := &PlanDebugger{}
	return VisitPlan(plan, debugger)
}

func (p *PlanDebugger) VisitCreateTable(plan *CreateTable) (string, error) {
	return "CreateTable: " + plan.Table.Name(), nil
}

func (p *PlanDebugger) VisitInsert(plan *Insert) (string, error) {
	return "Insert into: " + plan.Table.Name(), nil
}
