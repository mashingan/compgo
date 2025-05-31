package interp

import (
	"fmt"
)

type ObjectType string

const (
	IntegerType = "INTEGER"
	BooleanType = "BOOLEAN"
	NullType    = "NULL"
	RetType     = "RETURN"
	ErrorType   = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Primitive[T comparable] struct {
	Value T
}

func (i *Primitive[T]) Inspect() string { return fmt.Sprint(i.Value) }

type Integer struct {
	Primitive[int]
}

func (*Integer) Type() ObjectType { return IntegerType }

type Boolean struct {
	Primitive[bool]
}

func (*Boolean) Type() ObjectType { return BooleanType }

type Null struct {
	Primitive[*struct{}]
}

func (*Null) Type() ObjectType { return NullType }

type ReturnValue struct {
	Primitive[Object]
}

func (*ReturnValue) Type() ObjectType { return RetType }

type Error struct {
	Msg string
}

func (*Error) Type() ObjectType  { return ErrorType }
func (e *Error) Inspect() string { return fmt.Sprintf("ERROR: %s", e.Msg) }
