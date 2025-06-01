package interp

import (
	"fmt"
	"strings"
)

type ObjectType string

const (
	IntegerType    = "INTEGER"
	BooleanType    = "BOOLEAN"
	NullType       = "NULL"
	RetType        = "RETURN"
	ErrorType      = "ERROR"
	FunctionType   = "FUNCTION"
	IdentifierType = "IDENTIFIER"
	StringType     = "STRING"
	BuiltinType    = "BUILTIN"
	SliceType      = "ARRAY"
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

type Function struct {
	Parameters []*Identifier
	Body       *BlockStatement
	Env        *Environment
}

func (*Function) Type() ObjectType { return FunctionType }
func (f *Function) Inspect() string {
	prm := make([]string, len(f.Parameters))
	for i, s := range f.Parameters {
		prm[i] = s.String()
	}
	return fmt.Sprintf("fn (%s) {\n%s\n}", strings.Join(prm, ", "), f.Body)

}

type IdentifierObj struct {
	Primitive[string]
}

func (*IdentifierObj) Type() ObjectType { return IdentifierType }

type String struct {
	Primitive[string]
}

func (*String) Type() ObjectType  { return StringType }
func (s *String) Inspect() string { return fmt.Sprintf(`"%s"`, s.Value) }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (*Builtin) Type() ObjectType { return BuiltinType }
func (*Builtin) Inspect() string  { return "builtin-function" }

type SliceObj struct {
	Elements []Object
}

func (*SliceObj) Type() ObjectType { return SliceType }
func (s *SliceObj) Inspect() string {
	so := make([]string, len(s.Elements))
	for i, o := range s.Elements {
		so[i] = o.Inspect()
	}
	return fmt.Sprintf("[%s]", strings.Join(so, ","))
}
