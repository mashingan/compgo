package interp

import "testing"

func TestEvalInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"5", 5},
		{"10", 10},
	}
	for _, tt := range tests {
		ev := testEval(tt.input)
		testIntegerObject(t, ev, tt.expected)
	}
}

func testEval(input string) Object {
	p := NewParser(NewLexer(input))
	prg := p.ParseProgram()
	return Eval(prg.Statements[0])
}

func testIntegerObject(t *testing.T, o Object, expected int) bool {
	r, ok := o.(*Integer)
	if !ok {
		t.Errorf("%s is not integer. got=%T (%+v)", o.Inspect(), o, o)
		return false
	}
	if r.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			r.Value, expected)
		return false
	}
	return true
}
