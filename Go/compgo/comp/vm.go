package comp

import (
	"compgo/interp"
	"encoding/binary"
	"fmt"
)

const stackSize = 2048

type Vm struct {
	constants []interp.Object
	Instructions
	Stack[interp.Object]
	sp      int
	lastPop interp.Object
}

func NewVm(b *Bytecode) *Vm {
	return &Vm{
		Instructions: b.Instructions,
		constants:    b.Constants,
		Stack:        make(Stack[interp.Object], 0, stackSize),
		sp:           0,
	}
}

func (vm *Vm) StackTop() interp.Object {
	if len(vm.Stack) == 0 {
		return nil
	}
	return vm.Stack[len(vm.Stack)-1]
}

func (vm *Vm) Pop() (interp.Object, error) {
	v, err := vm.Stack.Pop()
	if err != nil {
		return nil, err
	}
	vm.lastPop = v
	return v, nil
}

func (vm *Vm) LastPop() interp.Object {
	return vm.lastPop
}

var (
	ErrEmptyStack = fmt.Errorf("empty stack")
)

type Stack[T any] []T

func (s *Stack[T]) Push(val T) {
	*s = append(*s, val)
}

func (s Stack[T]) Peek() (T, error) {
	if len(s) == 0 {
		var v T
		return v, ErrEmptyStack
	}
	return s[len(s)-1], nil
}

func (s *Stack[T]) Pop() (T, error) {
	if len(*s) == 0 {
		var val T
		return val, ErrEmptyStack
	}
	lens := len(*s) - 1
	val := (*s)[lens]
	*s = (*s)[:lens]
	return val, nil
}

func (vm *Vm) Run() error {
	ip := 0
	for ip < len(vm.Instructions) {
		op := Opcode(vm.Instructions[ip])
		ip++
		switch op {
		case OpConstant:
			idx := uint16(0)
			def := Definition{OperandWidth: []int{2}}
			binary.Decode(vm.Instructions[ip:ip+def.OperandWidth[0]],
				binary.BigEndian, &idx)
			vm.Stack.Push(vm.constants[idx])
			ip += def.OperandWidth[0]
		case OpAdd, OpSub, OpMul, OpDiv:
			fn, ok := mapInfixOps[op]
			if !ok {
				return fmt.Errorf("undefined infix operator: %d", op)
			}
			if err := fn(vm); err != nil {
				return err
			}
		case OpPop:
			vm.Pop()
		}
	}
	return nil
}

func (vm *Vm) pop2() (interp.Object, interp.Object, error) {
	left, err := vm.Pop()
	if err != nil {
		return nil, nil, err
	}
	right, err := vm.Pop()
	if err != nil {
		return nil, nil, err
	}
	return left, right, nil
}

var mapInfixOps = map[Opcode]func(vm *Vm) error{
	OpAdd: add,
	OpSub: sub,
	OpMul: mul,
	OpDiv: div,
}

func arith(vm *Vm, fop func(vm *Vm, left, right *interp.Integer)) error {
	lobj, robj, err := vm.pop2()
	if err != nil {
		return err
	}
	lint, ok := lobj.(*interp.Integer)
	if !ok {
		return fmt.Errorf("object is not integer. got=%T (%+v)", lobj, lobj)
	}
	rint, ok := robj.(*interp.Integer)
	if !ok {
		return fmt.Errorf("object is not integer. got=%T (%+v)", robj, robj)
	}
	fop(vm, lint, rint)
	return nil
}

func add(vm *Vm) error {
	return arith(vm, func(vm *Vm, left, right *interp.Integer) {
		left.Value += right.Value
		vm.Push(left)
	})
}

func sub(vm *Vm) error {
	return arith(vm, func(vm *Vm, left, right *interp.Integer) {
		left.Value -= right.Value
		vm.Push(left)
	})
}

func mul(vm *Vm) error {
	return arith(vm, func(vm *Vm, left, right *interp.Integer) {
		left.Value *= right.Value
		vm.Push(left)
	})
}

func div(vm *Vm) error {
	return arith(vm, func(vm *Vm, left, right *interp.Integer) {
		left.Value /= right.Value
		vm.Push(left)
	})
}
