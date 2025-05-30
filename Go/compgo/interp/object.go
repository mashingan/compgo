package interp

import (
	"fmt"
)

type ObjecType string

const (
	IntegerType = "INTEGER"
	BooleanType = "BOOLEAN"
	NullType    = "NULL"
)

type Object interface {
	Type() ObjecType
	Inspect() string
}

type Primitive[T comparable] struct {
	Value T
}

func (i *Primitive[T]) Inspect() string { return fmt.Sprint(i.Value) }

type Integer struct {
	Primitive[int]
}

func (*Integer) Type() ObjecType { return IntegerType }

type Boolean struct {
	Primitive[bool]
}

func (*Boolean) Type() ObjecType { return BooleanType }

type Null struct {
	Primitive[*struct{}]
}

func (*Null) Type() ObjecType { return NullType }
