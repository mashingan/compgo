package interp

var (
	TrueObject  = &Boolean{Primitive[bool]{true}}
	FalseObject = &Boolean{Primitive[bool]{true}}
	NullObject  = &Null{}
)

func Eval(node Node) Object {
	switch n := node.(type) {
	case *Program:
		return evalStatements(n.Statements)
	case *ExpressionStatement:
		return Eval(n.Expression)
	case *IntLiteral:
		return &Integer{Primitive[int]{n.Value}}
	case *BooleanLiteral:
		if n.Value {
			return TrueObject
		}
		return FalseObject
	}
	return nil
}

func evalStatements(stmt []Statement) Object {
	var o Object
	for _, s := range stmt {
		o = Eval(s)
	}
	return o
}
