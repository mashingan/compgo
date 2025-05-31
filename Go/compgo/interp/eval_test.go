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
	return Eval(prg)
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

func TestEvalBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, tt := range tests {
		ev := testEval(tt.input)
		testBooleanObject(t, ev, tt.expected)
	}
}

func testBooleanObject(t *testing.T, o Object, expected bool) bool {
	r, ok := o.(*Boolean)
	if !ok {
		t.Errorf("%s is not boolean. got=%T (%+v)", o.Inspect(), o, o)
		return false
	}
	if r.Value != expected {
		t.Errorf("'%s' has wrong value. got=%v, want=%v",
			o.Inspect(), r.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		t.Log("tt.input:", tt.input)
		evl := testEval(tt.input)
		testBooleanObject(t, evl, tt.expected)
	}
}
