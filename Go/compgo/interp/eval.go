package interp

import "fmt"

var (
	TrueObject  = &Boolean{Primitive[bool]{true}}
	FalseObject = &Boolean{Primitive[bool]{false}}
	NullObject  = &Null{}
)

const (
	unknownOperatorPrefixFmt = "unknown operator: %s%s"
	unknownOperatorInfixFmt  = "unknown operator: %s %s %s"
)

func Eval(node Node) Object {
	switch n := node.(type) {
	case *Program:
		return evalProgram(n.Statements)
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
		if _, yes := right.(*Error); yes {
			return right
		}
		return evalPrefix(n.Operator, right)
	case *InfixExpression:
		left := Eval(n.Left)
		if _, yes := left.(*Error); yes {
			return left
		}
		right := Eval(n.Right)
		if _, yes := right.(*Error); yes {
			return right
		}
		return evalInfix(n.Operator, left, right)
	case *BlockStatement:
		return evalBlockStatements(n.Statements)
	case *IfExpression:
		return evalIfElse(n)
	case *ReturnStatement:
		val := Eval(n.Value)
		if _, yes := val.(*Error); yes {
			return val
		}
		return &ReturnValue{Primitive[Object]{val}}
	}
	return nil
}

func evalProgram(stmt []Statement) Object {
	var o Object
	for _, s := range stmt {
		o = Eval(s)
		switch r := o.(type) {
		case *ReturnValue:
			return r.Value
		case *Error:
			return r
		}
	}
	return o
}

func evalBlockStatements(stmt []Statement) Object {
	var o Object
	for _, s := range stmt {
		o = Eval(s)
		if o != nil {
			rt := o.Type()
			if rt == RetType || rt == ErrorType {
				return o
			}
		}
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
			return &Error{fmt.Sprintf(unknownOperatorPrefixFmt, op, o.Type())}
		}
		i.Value *= -1
		return i
	default:
		return &Error{fmt.Sprintf(unknownOperatorPrefixFmt, op, o.Type())}
	}
}

func evalInfixMath(op string, left, right *Integer) Object {
	switch op {
	case "+":
		left.Value += right.Value
		return left
	case "-":
		left.Value -= right.Value
		return left
	case "*":
		left.Value *= right.Value
		return left
	case "/":
		left.Value /= right.Value
		return left
	default:
		return &Error{fmt.Sprintf(unknownOperatorInfixFmt,
			left.Type(), op, right.Type())}
	}
}

func evalInfix(op string, left, right Object) Object {
	lint, lok := left.(*Integer)
	rint, rok := right.(*Integer)
	switch op {
	case "+", "-", "*", "/":
		if !lok || !rok {
			return &Error{fmt.Sprintf(unknownOperatorInfixFmt,
				left.Type(), op, right.Type())}
		}
		return evalInfixMath(op, lint, rint)
	case "<=", ">=", ">", "<":
		if !lok || !rok {
			return NullObject
		}
		return evalCompareInt(op, lint, rint)
	case "==", "!=":
		if lok && rok {
			if op == "==" {
				if lint.Value == rint.Value {
					return TrueObject
				}
				return FalseObject
			}
			if lint.Value != rint.Value {
				return TrueObject
			}
			return FalseObject
		}
		b := toNativeBoolean(left) == toNativeBoolean(right)
		if (op == "==" && b) || (op == "!=" && !b) {
			return TrueObject
		}
		return FalseObject
	default:
		return &Error{fmt.Sprintf(unknownOperatorInfixFmt, left.Type(), op, right.Type())}
	}
}

func evalCompareInt(op string, left, right *Integer) Object {
	compare := func(test bool) Object {
		if test {
			return TrueObject
		}
		return FalseObject
	}
	switch op {
	case ">":
		return compare(left.Value > right.Value)
	case ">=":
		return compare(left.Value >= right.Value)
	case "<":
		return compare(left.Value < right.Value)
	case "<=":
		return compare(left.Value <= right.Value)
	default:
		return NullObject
	}
}

func toNativeBoolean(o Object) bool {
	switch vo := o.(type) {
	case *Integer:
		return vo.Value != 0
	case *Null:
		return false
	case *Boolean:
		return vo.Value
	default:
		return false
	}
}

func evalIfElse(ie *IfExpression) Object {
	cond := Eval(ie.Condition)
	if _, yes := cond.(*Error); yes {
		return cond
	}
	if toNativeBoolean(cond) {
		return Eval(ie.Then)
	} else if ie.Else != nil {
		return Eval(ie.Else)
	}
	return NullObject
}
