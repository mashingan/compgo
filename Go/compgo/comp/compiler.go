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

var mapOpCodes = map[string]Opcode{
	"+":  OpAdd,
	"-":  OpSub,
	"*":  OpMul,
	"/":  OpDiv,
	"==": OpEq,
	"!=": OpNeq,
	">":  OpGt,
	"<":  OpLt,
	">=": OpGte,
	"<=": OpLte,
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
		c.emit(OpPop)
	case *interp.InfixExpression:
		err := c.Compile(n.Left)
		if err != nil {
			return err
		}
		err = c.Compile(n.Right)
		if err != nil {
			return err
		}

		nop, ok := mapOpCodes[n.Operator]
		if !ok {
			return fmt.Errorf("unknown operator %s", n.Operator)
		}
		c.emit(nop)
	case *interp.PrefixExpression:
		err := c.Compile(n.Right)
		if err != nil {
			return nil
		}
		switch n.Operator {
		case "-":
			c.emit(OpMinus)
		case "!":
			c.emit(OpBang)
		default:
			return fmt.Errorf("unknown operator %s", n.Operator)
		}
	case *interp.IntLiteral:
		itg := &interp.Integer{Primitive: interp.Primitive[int]{Value: n.Value}}
		c.constants = append(c.constants, itg)
		c.emit(OpConstant, len(c.constants)-1)
	case *interp.BooleanLiteral:
		if n.Value {
			c.emit(OpTrue)
		} else {
			c.emit(OpFalse)
		}
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
