package comp

import (
	"encoding/binary"
	"fmt"
)

type (
	Instructions []byte
	Opcode       byte
)

const (
	OpConstant Opcode = iota
)

type Definition struct {
	Name         string
	OperandWidth []int
}

var definitions = map[Opcode]Definition{
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op byte) (Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return Definition{}, fmt.Errorf("op '%d' is undefined", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return nil
	}
	instlen := 1
	for _, w := range def.OperandWidth {
		instlen += w
	}
	inst := make([]byte, instlen)
	inst[0] = byte(op)
	offset := 1
	for i, o := range operands {
		width := def.OperandWidth[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(inst[offset:], uint16(o))
		}
		offset += width
	}
	return inst
}
