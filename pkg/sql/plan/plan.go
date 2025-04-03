package plan

type PlanVisitor[T any] interface {
	VisitProjection(*Projection) (*T, error)
	VisitInsert(*Insert) (*T, error)
}

type Plan interface {
	isPlan()
}

type Projection struct{}

func (p *Projection) isPlan() {}

type Insert struct {
	table  string
	source Plan
}

func (p *Insert) isPlan() {}

type Rows struct{}

func (p *Rows) isPlan() {}
