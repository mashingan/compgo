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
	Name string
	fn   *interp.Builtin
}{
	{
		pos:  0,
		Name: "len",
		fn:   interp.Builtins["len"],
	},
	{
		pos:  1,
		Name: "first",
		fn:   interp.Builtins["first"],
	},
	{
		pos:  2,
		Name: "last",
		fn:   interp.Builtins["last"],
	},
	{
		pos:  3,
		Name: "rest",
		fn:   interp.Builtins["rest"],
	},
	{
		pos:  4,
		Name: "push",
		fn:   interp.Builtins["push"],
	},
	{
		pos:  5,
		Name: "puts",
		fn:   interp.Builtins["puts"],
	},
}
