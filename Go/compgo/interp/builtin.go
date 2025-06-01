package interp

import (
	"fmt"
	"unicode/utf8"
)

var builtins = map[string]*Builtin{
	"len": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *String:
				return &Integer{Primitive[int]{utf8.RuneCountInString(arg.Value)}}
			default:
				return &Error{fmt.Sprintf("argument to 'len' not supported, got %s",
					args[0].Type())}
			}
		},
	},
}
