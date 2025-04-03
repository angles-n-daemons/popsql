package plan

type PlanDebugger struct{}

func (p *PlanDebugger) Debug(plan Plan) string {
	switch plan := plan.(type) {
	case *CreateTable:
		return "CreateTable: " + plan.Table.Name()
	default:
		return "Unknown Plan"
	}
}
