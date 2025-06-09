package comp

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpSub, []int{}, []byte{byte(OpSub)}},
		{OpMul, []int{}, []byte{byte(OpMul)}},
		{OpDiv, []int{}, []byte{byte(OpDiv)}},
		{OpEq, []int{}, []byte{byte(OpEq)}},
		{OpNeq, []int{}, []byte{byte(OpNeq)}},
		{OpGt, []int{}, []byte{byte(OpGt)}},
		{OpLt, []int{}, []byte{byte(OpLt)}},
		{OpGte, []int{}, []byte{byte(OpGte)}},
		{OpLte, []int{}, []byte{byte(OpLte)}},
		{OpPop, []int{}, []byte{byte(OpPop)}},
		{OpTrue, []int{}, []byte{byte(OpTrue)}},
		{OpFalse, []int{}, []byte{byte(OpFalse)}},
		{OpBang, []int{}, []byte{byte(OpBang)}},
		{OpMinus, []int{}, []byte{byte(OpMinus)}},
		{OpJump, []int{12}, []byte{byte(OpJump), 0, 12}},
		{OpJumpIfFalsy, []int{12}, []byte{byte(OpJumpIfFalsy), 0, 12}},
		{OpNull, []int{}, []byte{byte(OpNull)}},
		{OpGetLocal, []int{255}, []byte{byte(OpGetLocal), 255}},
		{OpSetLocal, []int{255}, []byte{byte(OpSetLocal), 255}},
		{OpCall, []int{255}, []byte{byte(OpCall), 255}},
	}
	for _, tt := range tests {
		inst := Make(tt.op, tt.operands...)
		if len(inst) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want=%d got=%d",
				len(tt.expected), len(inst))
		}
		for i, b := range tt.expected {
			if inst[i] != b {
				t.Errorf("wrong byte at pos %d. want=%d got=%d", i, b, inst[i])
			}
		}
	}
}
