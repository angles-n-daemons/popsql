package tree

type Visitor[T any] interface {
	// VisitBinary
}

// if i just think for a second, how do these things lend themselves to being walked.
// Let's say I wanted to do three different things along the nodes of the tree
// Interpret -> any
// Debug -> err
// SearchForString ->
type Expr[T any] interface {
	Accept(visitor Visitor[T])
	Walk(func(Expr[T]) error) error
	Type() string
}

// If each node had a walk function, it wouldn't work for interpreting.
// walk would have to still need to assert a type
// Interpreter.walk() {
//    switch type(node) {
//       case Binary {
//          left = left.accept(me)
//          return ...
// But if I added a new type, there would be no knowing if it were missed by the interpreter.
// Seems kinda safer, though now the issue becomes returning the right value
