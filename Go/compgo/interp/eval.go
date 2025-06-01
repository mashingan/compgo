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

func Eval(node Node, env *Environment) Object {
	switch n := node.(type) {
	case *Program:
		return evalProgram(n.Statements, env)
	case *ExpressionStatement:
		return Eval(n.Expression, env)
	case *IntLiteral:
		return &Integer{Primitive[int]{n.Value}}
	case *BooleanLiteral:
		if n.Value {
			return TrueObject
		}
		return FalseObject
	case *PrefixExpression:
		right := Eval(n.Right, env)
		if _, yes := right.(*Error); yes {
			return right
		}
		return evalPrefix(n.Operator, right)
	case *InfixExpression:
		left := Eval(n.Left, env)
		if _, yes := left.(*Error); yes {
			return left
		}
		right := Eval(n.Right, env)
		if _, yes := right.(*Error); yes {
			return right
		}
		return evalInfix(n.Operator, left, right)
	case *BlockStatement:
		return evalBlockStatements(n.Statements, env)
	case *IfExpression:
		return evalIfElse(n, env)
	case *ReturnStatement:
		val := Eval(n.Value, env)
		if _, yes := val.(*Error); yes {
			return val
		}
		return &ReturnValue{Primitive[Object]{val}}
	case *LetStatement:
		val := Eval(n.Value, env)
		if _, yes := val.(*Error); yes {
			return val
		}
		env.Set(n.Name.Value, val)
	case *Identifier:
		return evalIdentifier(n, env)
	case *FuncLiteral:
		params := n.Parameters
		body := n.Body
		return &Function{Parameters: params, Env: env, Body: body}
	case *CallExpression:
		fn := Eval(n.Func, env)
		if _, yes := fn.(*Error); yes {
			return fn
		}
		args := evalExpression(n.Args, env)
		if len(args) == 1 {
			if _, yes := args[0].(*Error); yes {
				return args[0]
			}
		}
		return evalCall(fn, args)
	}
	return nil
}

func evalProgram(stmt []Statement, env *Environment) Object {
	var o Object
	for _, s := range stmt {
		o = Eval(s, env)
		switch r := o.(type) {
		case *ReturnValue:
			return r.Value
		case *Error:
			return r
		}
	}
	return o
}

func evalBlockStatements(stmt []Statement, env *Environment) Object {
	var o Object
	for _, s := range stmt {
		o = Eval(s, env)
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

func evalIfElse(ie *IfExpression, env *Environment) Object {
	cond := Eval(ie.Condition, env)
	if _, yes := cond.(*Error); yes {
		return cond
	}
	if toNativeBoolean(cond) {
		return Eval(ie.Then, env)
	} else if ie.Else != nil {
		return Eval(ie.Else, env)
	}
	return NullObject
}

func evalIdentifier(o *Identifier, env *Environment) Object {
	val, ok := env.Get(o.Value)
	if !ok {
		return &Error{fmt.Sprintf("identifier not found: %s", o.Value)}
	}
	return val
}

func evalExpression(exps []Expression, env *Environment) []Object {
	res := make([]Object, len(exps))
	for i, e := range exps {
		evl := Eval(e, env)
		if _, yes := evl.(*Error); yes {
			return []Object{evl}
		}
		res[i] = evl
	}
	return res
}

func evalCall(fn Object, args []Object) Object {
	ffn, ok := fn.(*Function)
	if !ok {
		return &Error{fmt.Sprintf("not a function: %s", fn.Type())}
	}
	envFrame := NewEnvironmentFrame(ffn.Env)
	for i, a := range ffn.Parameters {
		envFrame.Set(a.Value, args[i])
	}
	evl := Eval(ffn.Body, envFrame)
	if val, ok := evl.(*ReturnValue); ok {
		return val.Value
	}
	return evl
}
