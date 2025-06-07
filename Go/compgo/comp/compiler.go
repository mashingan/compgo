package comp

import "compgo/interp"

type Compiler struct {
	Instructions
	constants []interp.Object
}

func New() *Compiler {
	return &Compiler{
		Instructions: Instructions{},
		constants:    []interp.Object{},
	}
}

func (c *Compiler) Compile(node interp.Node) error {
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.Instructions,
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions
	Constants []interp.Object
}
