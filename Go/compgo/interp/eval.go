package interp

var (
	TrueObject  = &Boolean{Primitive[bool]{true}}
	FalseObject = &Boolean{Primitive[bool]{false}}
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
	case *PrefixExpression:
		right := Eval(n.Right)
		return evalPrefix(n.Operator, right)
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

func evalPrefix(op string, o Object) Object {
	switch op {
	case "!":
		switch o {
		case TrueObject:
			return FalseObject
		case FalseObject:
			return TrueObject
		case NullObject:
			return TrueObject
		default:
			return FalseObject
		}
	case "-":
		i, ok := o.(*Integer)
		if !ok {
			return NullObject
		}
		i.Value *= -1
		return i
	default:
		return NullObject
	}
}
