package comp

import (
	"compgo/interp"
	"fmt"
)

const (
	CompiledFuncType interp.ObjectType = "COMPILED_FUNCTION_OBJ"
)

type CompiledFunction struct {
	Instructions
	NumLocals int
	NumArgs   int
}

func (c *CompiledFunction) Type() interp.ObjectType { return CompiledFuncType }
func (c *CompiledFunction) Inspect() string {
	return fmt.Sprintf("Compiledfunction[%p]", c)
}

var builtins = map[string]struct {
	pos int
	fn  *interp.Builtin
}{
	"len": {
		pos: 0,
		fn:  interp.Builtins["len"],
	},
	"first": {
		pos: 1,
		fn:  interp.Builtins["first"],
	},
	"last": {
		pos: 2,
		fn:  interp.Builtins["last"],
	},
	"rest": {
		pos: 3,
		fn:  interp.Builtins["rest"],
	},
	"push": {
		pos: 4,
		fn:  interp.Builtins["push"],
	},
	"puts": {
		pos: 5,
		fn:  interp.Builtins["puts"],
	},
}
