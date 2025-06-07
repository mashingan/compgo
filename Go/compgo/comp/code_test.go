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
