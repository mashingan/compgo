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
		case OpAdd:
			left, err := vm.Pop()
			if err != nil {
				return err
			}
			right, err := vm.Pop()
			if err != nil {
				return err
			}
			lint, ok := left.(*interp.Integer)
			if !ok {
				return fmt.Errorf("object is not integer. got=%T (%+v)", left, left)
			}
			rint, ok := right.(*interp.Integer)
			if !ok {
				return fmt.Errorf("object is not integer. got=%T (%+v)", right, right)
			}
			lint.Value += rint.Value
			vm.Push(lint)
		case OpPop:
			vm.Pop()
		}
	}
	return nil
}
