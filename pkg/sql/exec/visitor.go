package exec

type Record interface {
	Get(...string)
}

type Executor struct{}

func (e *Executor) VisitCreateTableStmt() ([]*Record, error) {
	return nil, nil
}

func (e *Executor) VisitSelectStmt() ([]*Record, error) {
	return nil, nil
}

func (e *Executor) VisitInsertStmt() ([]*Record, error) {
	return nil, nil
}
