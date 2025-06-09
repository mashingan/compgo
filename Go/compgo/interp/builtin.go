package interp

import (
	"fmt"
	"unicode/utf8"
)

func wrongArguments(expected, got int) Object {
	return &Error{fmt.Sprintf("wrong number of arguments. got=%d, want=%d",
		got, expected)}
}

var Builtins = map[string]*Builtin{
	"len": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return wrongArguments(1, len(args))
			}
			switch arg := args[0].(type) {
			case *String:
				return &Integer{Primitive[int]{utf8.RuneCountInString(arg.Value)}}
			case *SliceObj:
				return &Integer{Primitive[int]{len(arg.Elements)}}
			default:
				return &Error{fmt.Sprintf("argument to 'len' not supported, got %s",
					args[0].Type())}
			}
		},
	},
	"first": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return wrongArguments(1, len(args))
			}
			switch arg := args[0].(type) {
			case *String:
				var r rune = 0
				for _, rr := range arg.Value {
					r = rr
					break
				}
				if r == 0 {
					return NullObject
				}
				return &String{Primitive[string]{string(r)}}
			case *SliceObj:
				if len(arg.Elements) < 1 {
					return NullObject
				}
				return arg.Elements[0]
			default:
				return &Error{fmt.Sprintf("argument to 'first' not supported, got %s",
					args[0].Type())}
			}
		},
	},
	"last": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return wrongArguments(1, len(args))
			}
			switch arg := args[len(args)-1].(type) {
			case *String:
				r, sz := utf8.DecodeLastRuneInString(arg.Value)
				if sz == 0 {
					return NullObject
				}
				return &String{Primitive[string]{string(r)}}
			case *SliceObj:
				if len(arg.Elements) < 1 {
					return NullObject
				}
				return arg.Elements[len(arg.Elements)-1]
			default:
				return &Error{fmt.Sprintf("argument to 'last' not supported, got %s",
					args[len(args)-1].Type())}
			}
		},
	},
	"rest": {
		func(args ...Object) Object {
			if len(args) != 1 {
				return wrongArguments(1, len(args))
			}
			switch arg := args[0].(type) {
			case *String:
				length := utf8.RuneCountInString(arg.Value)
				if length < 1 {
					return NullObject
				}
				_, sz := utf8.DecodeRuneInString(arg.Value)
				arg.Value = arg.Value[sz:]
				return arg
			case *SliceObj:
				length := len(arg.Elements)
				if length < 1 {
					return NullObject
				}
				arg.Elements = arg.Elements[1:]
				return arg
			default:
				return &Error{fmt.Sprintf("argument to 'rest' not supported, got %s",
					args[len(args)-1].Type())}

			}
		},
	},
	"push": {
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return wrongArguments(2, 0)
			} else if len(args) == 1 {
				return args[0]
			}
			switch arg := args[0].(type) {
			case *String:
				for _, s := range args[1:] {
					switch ss := s.(type) {
					case *String:
						arg.Value += ss.Value
					case *Integer:
						arg.Value += fmt.Sprint(arg.Value)
					case *Boolean:
						arg.Value += fmt.Sprint(arg.Value)
					default:
						arg.Value += s.Inspect()
					}
				}
				return arg
			case *SliceObj:
				arg.Elements = append(arg.Elements, args[1:]...)
				return arg
			default:
				return &Error{fmt.Sprintf("argument to 'push' not supported, got %s",
					args[0].Type())}
			}
		},
	},
	"puts": {
		Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NullObject
		},
	},
}
