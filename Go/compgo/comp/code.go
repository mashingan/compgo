package comp

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type (
	Instructions []byte
	Opcode       byte
)

const (
	OpConstant Opcode = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpNot
	OpTrue
	OpFalse
	OpEq
	OpNeq
	OpGt
	OpLt
	OpGte
	OpLte
	OpPop
	OpBang
	OpMinus
	OpJumpIfFalsy
	OpJump
	OpNull
	OpGetGlobal
	OpSetGlobal
	OpArray
	OpHash
	OpIndex
	OpCall
	OpReturn
	OpReturnValue
	OpGetLocal
	OpSetLocal
	OpGetBuiltin
	OpClosure
)

type Definition struct {
	Name         string
	OperandWidth []int
}

var definitions = map[Opcode]Definition{
	OpConstant:    {"OpConstant", []int{2}},
	OpAdd:         {"OpAdd", []int{}},
	OpSub:         {"OpSub", []int{}},
	OpMul:         {"OpMul", []int{}},
	OpDiv:         {"OpDiv", []int{}},
	OpEq:          {"OpEq", []int{}},
	OpNeq:         {"OpNeq", []int{}},
	OpGt:          {"OpGt", []int{}},
	OpLt:          {"OpLt", []int{}},
	OpGte:         {"OpGte", []int{}},
	OpLte:         {"OpLte", []int{}},
	OpPop:         {"OpPop", []int{}},
	OpTrue:        {"OpTrue", []int{}},
	OpFalse:       {"OpFalse", []int{}},
	OpBang:        {"OpBang", []int{}},
	OpMinus:       {"OpMinus", []int{}},
	OpJump:        {"OpJump", []int{2}},
	OpJumpIfFalsy: {"OpJumpIfFalsy", []int{2}},
	OpNull:        {"OpNull", []int{}},
	OpGetGlobal:   {"OpGetGlobal", []int{2}},
	OpSetGlobal:   {"OpSetGlobal", []int{2}},
	OpArray:       {"OpArray", []int{2}},
	OpHash:        {"OpHash", []int{2}},
	OpIndex:       {"OpIndex", []int{}},
	OpCall:        {"OpCall", []int{1}},
	OpReturn:      {"OpReturn", []int{}},
	OpReturnValue: {"OpReturnValue", []int{}},
	OpGetLocal:    {"OpGetLocal", []int{1}},
	OpSetLocal:    {"OpSetLocal", []int{1}},
	OpGetBuiltin:  {"OpGetBuiltin", []int{1}},
	OpClosure:     {"OpClosure", []int{2, 1}},
}

func (i Instructions) String() string {
	var (
		sb   strings.Builder
		addr int
	)
	for addr < len(i) {
		sb.WriteString(fmt.Sprintf("%04d ", addr))
		def, ok := definitions[Opcode(i[addr])]
		addr++
		if !ok {
			sb.WriteByte('\n')
			continue
		}
		sb.WriteString(def.Name)
		for _, lond := range def.OperandWidth {
			switch lond {
			case 1:
				opr := i[addr]
				addr += lond
				sb.WriteString(fmt.Sprintf(" %d", opr))
			case 2:
				val := binary.BigEndian.Uint16(i[addr : addr+lond])
				addr += lond
				sb.WriteString(fmt.Sprintf(" %d", val))

			}
		}
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
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
		case 1:
			inst[offset] = byte(o)
		case 2:
			binary.BigEndian.PutUint16(inst[offset:], uint16(o))
		}
		offset += width
	}
	return inst
}
