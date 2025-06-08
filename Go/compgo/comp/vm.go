package comp

import (
	"compgo/interp"
	"encoding/binary"
	"errors"
	"fmt"
	"path"
	"runtime"
	"unicode/utf8"
)

const (
	stackSize  = 2048
	GlobalSize = 65536
)

type Vm struct {
	constants []interp.Object
	Instructions
	Stack[interp.Object]
	sp      int
	lastPop interp.Object
	globals []interp.Object
}

func NewVm(b *Bytecode) *Vm {
	return &Vm{
		Instructions: b.Instructions,
		constants:    b.Constants,
		Stack:        make(Stack[interp.Object], 0, stackSize),
		sp:           0,
		globals:      make([]interp.Object, GlobalSize),
	}
}

func (vm *Vm) SetConstants(cnts []interp.Object) {
	vm.constants = cnts
}

func (vm *Vm) SetGlobals(globs []interp.Object) {
	vm.globals = globs
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
	inspectEmptyStack := func(err error) {
		pc, file, lineno, _ := runtime.Caller(1)
		funcname := runtime.FuncForPC(pc).Name()
		fname := path.Base(file)
		if errors.Is(err, ErrEmptyStack) {
			fmt.Printf("%s#%s:%d inst:\n%s\n",
				fname, funcname, lineno, vm.Instructions)
		}
	}
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
		case OpAdd, OpSub, OpMul, OpDiv, OpEq, OpNeq,
			OpLt, OpLte, OpGt, OpGte:
			fn, ok := mapInfixOps[op]
			if !ok {
				return fmt.Errorf("undefined infix operator: %d", op)
			}
			if err := fn(vm); err != nil {
				inspectEmptyStack(err)
				return err
			}
		case OpPop:
			vm.Pop()
		case OpTrue:
			vm.Push(interp.TrueObject)
		case OpFalse:
			vm.Push(interp.FalseObject)
		case OpMinus:
			lastval, err := vm.Pop()
			if err != nil {
				inspectEmptyStack(err)
				return err
			}
			i, ok := lastval.(*interp.Integer)
			if !ok {
				return fmt.Errorf("wrong type not integer. got=%T (%+v)",
					lastval, lastval)
			}
			i.Value *= -1
			vm.Push(i)
		case OpBang:
			lastitem, err := vm.Pop()
			if err != nil {
				inspectEmptyStack(err)
				return err
			}
			if err := notObj(vm, lastitem); err != nil {
				inspectEmptyStack(err)
				return err
			}
		case OpJump:
			addr := uint16(0)
			binary.Decode(vm.Instructions[ip:], binary.BigEndian, &addr)
			ip = int(addr)
		case OpJumpIfFalsy:
			addr := uint16(0)
			binary.Decode(vm.Instructions[ip:], binary.BigEndian, &addr)
			ip += 2
			cond, err := vm.Pop()
			if err != nil {
				inspectEmptyStack(err)
				return err
			}
			if !isTruthy(cond) {
				ip = int(addr)
			}
		case OpNull:
			vm.Push(interp.NullObject)
		case OpSetGlobal:
			idx := uint16(0)
			binary.Decode(vm.Instructions[ip:], binary.BigEndian, &idx)
			ip += 2
			glb, err := vm.Pop()
			if err != nil {
				inspectEmptyStack(err)
				return err
			}
			vm.globals[idx] = glb
		case OpGetGlobal:
			idx := uint16(0)
			binary.Decode(vm.Instructions[ip:], binary.BigEndian, &idx)
			ip += 2
			glb := vm.globals[idx]
			vm.Push(glb)
		case OpArray:
			elm := uint16(0)
			binary.Decode(vm.Instructions[ip:], binary.BigEndian, &elm)
			ip += 2
			vm.sp = len(vm.Stack) - int(elm)
			arr := &interp.SliceObj{Elements: make([]interp.Object, elm)}
			for i := vm.sp; i < len(vm.Stack); i++ {
				arr.Elements[i-vm.sp] = vm.Stack[i]
			}
			vm.Push(arr)
		case OpHash:
			pairs := uint16(0)
			binary.Decode(vm.Instructions[ip:], binary.BigEndian, &pairs)
			ip += 2
			vm.sp = len(vm.Stack) - int(pairs)
			h := &interp.Hash{Pairs: map[interp.HashKey]interp.HashPair{}}
			for i := vm.sp; i < len(vm.Stack); i += 2 {
				k := vm.Stack[i]
				v := vm.Stack[i+1]
				pair := interp.HashPair{Key: k, Value: v}
				hk, ok := k.(interp.Hashable)
				if !ok {
					return fmt.Errorf("unusable as hash key: %s", k.Type())
				}
				h.Pairs[hk.HashKey()] = pair
			}
			vm.Push(h)
		case OpIndex:
			if err := processIndex(vm); err != nil {
				inspectEmptyStack(err)
				return err
			}
		}
	}
	return nil
}

func (vm *Vm) pop2() (interp.Object, interp.Object, error) {
	right, err := vm.Pop()
	if err != nil {
		return nil, nil, err
	}
	left, err := vm.Pop()
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
	OpEq:  eqObj,
	OpNeq: neqObj,
	OpGt:  gtObj,
	OpLt:  ltObj,
	OpGte: gteObj,
	OpLte: lteObj,
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
	lobj, robj, err := vm.pop2()
	if err != nil {
		return err
	}
	switch rint := robj.(type) {
	case *interp.Integer:
		lint, ok := lobj.(*interp.Integer)
		if !ok {
			return fmt.Errorf("unknown operator: %s + %s", lobj.Type(), robj.Type())
		}
		newv := &interp.Integer{Primitive: interp.Primitive[int]{
			Value: lint.Value + rint.Value,
		}}
		vm.Push(newv)
	case *interp.String:
		lstr, ok := lobj.(*interp.String)
		if !ok {
			return fmt.Errorf("unknown operator: %s + %s", lobj.Type(), robj.Type())
		}
		newv := &interp.String{Primitive: interp.Primitive[string]{
			Value: lstr.Value + rint.Value,
		}}
		vm.Push(newv)
	case *interp.SliceObj:
		var (
			larr  *interp.SliceObj
			lleft interp.Object
			ok    bool
		)
		pos := len(vm.Stack) - 1
		for i := pos; i >= 0; i-- {
			lleft = vm.Stack[i]
			larr, ok = lleft.(*interp.SliceObj)
			if ok {
				break
			}
			pos--
		}
		if pos == 0 {
			return fmt.Errorf("unknown operator: %s + %s", lleft.Type(), robj.Type())
		}
		for range len(vm.Stack) - pos {
			vm.Pop()
		}
		newarr := &interp.SliceObj{Elements: []interp.Object{}}
		newarr.Elements = append(newarr.Elements, larr.Elements...)
		newarr.Elements = append(newarr.Elements, rint.Elements...)
		vm.Push(newarr)
	default:
		return fmt.Errorf("unknown operator: %s + %s", lobj.Type(), robj.Type())
	}
	return nil
}

func sub(vm *Vm) error {
	return arith(vm, func(vm *Vm, left, right *interp.Integer) {
		newv := &interp.Integer{Primitive: interp.Primitive[int]{
			Value: left.Value - right.Value,
		}}
		vm.Push(newv)
	})
}

func mul(vm *Vm) error {
	return arith(vm, func(vm *Vm, left, right *interp.Integer) {
		newv := &interp.Integer{Primitive: interp.Primitive[int]{
			Value: left.Value * right.Value,
		}}
		vm.Push(newv)
	})
}

func div(vm *Vm) error {
	return arith(vm, func(vm *Vm, left, right *interp.Integer) {
		newv := &interp.Integer{Primitive: interp.Primitive[int]{
			Value: left.Value / right.Value,
		}}
		vm.Push(newv)
	})
}

func comparableObj(vm *Vm, test func(l, r interp.Object) bool) error {
	left, right, err := vm.pop2()
	if err != nil {
		return err
	}
	if left.Type() != right.Type() {
		return fmt.Errorf("not the same object type, left=%q and right=%q",
			left.Inspect(), right.Inspect())
	}
	if test(left, right) {
		vm.Push(interp.TrueObject)
	} else {
		vm.Push(interp.FalseObject)
	}
	return nil
}

func eqObj(vm *Vm) error {
	return comparableObj(vm, func(l, r interp.Object) bool {
		switch left := l.(type) {
		case *interp.Integer:
			right := r.(*interp.Integer)
			return left.Value == right.Value
		case *interp.Boolean:
			right := r.(*interp.Boolean)
			return left.Value == right.Value
		}
		return false
	})
}

func notObj(vm *Vm, obj interp.Object) error {
	switch b := obj.(type) {
	case *interp.Boolean:
		if b.Value {
			vm.Push(interp.FalseObject)
		} else {
			vm.Push(interp.TrueObject)
		}
		return nil
	case *interp.Integer:
		if b.Value == 0 {
			vm.Push(interp.TrueObject)
		} else {
			vm.Push(interp.FalseObject)
		}
		return nil
	case *interp.Null:
		vm.Push(interp.TrueObject)
		return nil
	default:
		return fmt.Errorf("cannot be applied for not-equality. got=%T (%+v)",
			obj, obj)
	}
}

func neqObj(vm *Vm) error {
	if err := eqObj(vm); err != nil {
		return err
	}
	lastBool, err := vm.Pop()
	if err != nil {
		return err
	}
	return notObj(vm, lastBool)
}

func orderableObj(vm *Vm, test func(l, r *interp.Integer) bool) error {
	left, right, err := vm.pop2()
	if err != nil {
		return err
	}
	if left.Type() != right.Type() {
		return fmt.Errorf("not the same object type, left=%q and right=%q",
			left.Inspect(), right.Inspect())
	}
	switch lobj := left.(type) {
	case *interp.Integer:
		robj := right.(*interp.Integer)
		if test(lobj, robj) {
			vm.Push(interp.TrueObject)
			return nil
		}
		vm.Push(interp.FalseObject)
	}
	return nil
}

func gtObj(vm *Vm) error {
	return orderableObj(vm, func(l, r *interp.Integer) bool {
		return l.Value > r.Value
	})
}

func ltObj(vm *Vm) error {
	return orderableObj(vm, func(l, r *interp.Integer) bool {
		return l.Value < r.Value
	})
}

func gteObj(vm *Vm) error {
	return orderableObj(vm, func(l, r *interp.Integer) bool {
		return l.Value >= r.Value
	})
}

func lteObj(vm *Vm) error {
	return orderableObj(vm, func(l, r *interp.Integer) bool {
		return l.Value <= r.Value
	})
}

func isTruthy(o interp.Object) bool {
	switch b := o.(type) {
	case *interp.Boolean:
		return b.Value
	case *interp.Integer:
		return b.Value != 0
	case *interp.Null:
		return false
	default:
		return true
	}
}

func processIndex(vm *Vm) error {
	left, idx, err := vm.pop2()
	if err != nil {
		return err
	}
	switch lobj := left.(type) {
	case *interp.SliceObj:
		idn, ok := idx.(*interp.Integer)
		if !ok {
			return fmt.Errorf("index accessing array is not integer. got=%T (%+v)",
				idx, idx)
		}
		if idn.Value >= len(lobj.Elements) || idn.Value < 0 {
			vm.Push(interp.NullObject)
			return nil
		}
		vm.Push(lobj.Elements[idn.Value])
	case *interp.String:
		idn, ok := idx.(*interp.Integer)
		if !ok {
			return fmt.Errorf("index accessing string is not integer. got=%T (%+v)",
				idx, idx)
		}
		strlen := utf8.RuneCountInString(lobj.Value)
		if idn.Value >= strlen || idn.Value < 0 {
			vm.Push(interp.NullObject)
			return nil
		}
		str := &interp.String{Primitive: interp.Primitive[string]{
			Value: "",
		}}
		count := 0
		for _, s := range lobj.Value {
			if count == idn.Value {
				str.Value = string(s)
				break
			}
			count++
		}
		vm.Push(str)
	case *interp.Hash:
		h, ok := idx.(interp.Hashable)
		if !ok {
			return fmt.Errorf("unusable key as hash: %s", idx.Inspect())
		}
		o, ok := lobj.Pairs[h.HashKey()]
		if !ok {
			vm.Push(interp.NullObject)
			return nil
		}
		vm.Push(o.Value)
	}
	return nil
}
