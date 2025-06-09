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
}

func (c *CompiledFunction) Type() interp.ObjectType { return CompiledFuncType }
func (c *CompiledFunction) Inspect() string {
	return fmt.Sprintf("Compiledfunction[%p]", c)
}
