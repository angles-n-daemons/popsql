package grammar

type walkFunc func(Expr[T]) error

type Expr[T any] interface {
	Accept(visitor Visitor[T]) error
	Walk(walkFunc) error
}

type ExprVisitor interface {
	VisitBinaryExpr(Binary) error
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (e Binary) Walk(f walkFunc) error {
	var err error
	err = e.left.Walk(f)

	if err != nil {
		return err
	}

	err = e.operator.Walk(f)

	if err != nil {
		return err
	}

	err = e.right.Walk(f)

	if err != nil {
		return err
	}

	return nil
}

func (e Binary) Visit(v ExprVisitor) error {
	return v.VisitBinaryExpr(e)
}
