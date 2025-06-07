package comp

import (
	"compgo/interp"
	"fmt"
)

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
	switch n := node.(type) {
	case *interp.Program:
		for _, s := range n.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *interp.ExpressionStatement:
		err := c.Compile(n.Expression)
		if err != nil {
			return err
		}
	case *interp.InfixExpression:
		err := c.Compile(n.Left)
		if err != nil {
			return err
		}
		err = c.Compile(n.Right)
		if err != nil {
			return err
		}
		switch n.Operator {
		case "+":
			c.emit(OpAdd)
		default:
			return fmt.Errorf("unknown operator %s", n.Operator)
		}
	case *interp.IntLiteral:
		itg := &interp.Integer{Primitive: interp.Primitive[int]{Value: n.Value}}
		c.constants = append(c.constants, itg)
		c.emit(OpConstant, len(c.constants)-1)
	}
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.Instructions,
		Constants:    c.constants,
	}
}

func (c *Compiler) emit(op Opcode, operands ...int) int {
	ins := Make(op, operands...)
	pos := len(c.Instructions)
	c.Instructions = append(c.Instructions, ins...)
	// return len(c.Instructions) - 1
	return pos
}

type Bytecode struct {
	Instructions
	Constants []interp.Object
}
