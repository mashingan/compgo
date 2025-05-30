package interp

func Eval(node Node) Object {
	switch n := node.(type) {
	case *Program:
		return evalStatements(n.Statements)
	case *ExpressionStatement:
		return Eval(n.Expression)
	case *IntLiteral:
		return &Integer{Primitive[int]{n.Value}}
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
