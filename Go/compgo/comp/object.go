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

var Builtins = []struct {
	pos  int
	name string
	fn   *interp.Builtin
}{
	{
		pos:  0,
		name: "len",
		fn:   interp.Builtins["len"],
	},
	{
		pos:  1,
		name: "first",
		fn:   interp.Builtins["first"],
	},
	{
		pos:  2,
		name: "last",
		fn:   interp.Builtins["last"],
	},
	{
		pos:  3,
		name: "rest",
		fn:   interp.Builtins["rest"],
	},
	{
		pos:  4,
		name: "push",
		fn:   interp.Builtins["push"],
	},
	{
		pos:  5,
		name: "puts",
		fn:   interp.Builtins["puts"],
	},
}
