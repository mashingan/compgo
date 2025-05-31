package interp

import "testing"

func TestEvalInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"--5", 5},
		{"--10", 10},
		{"5+5+5+5-10", 10},
		{"2*2*2*2*2", 32},
		{"-50 + 100 +-50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 *-10", 0},
		{"50 * 2/2 +10", 60},
		{"2 * (5+10)", 30},
		{"3*3*3+ 10", 37},
		{"3*(3*3)+ 10", 37},
		{"( 5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		ev := testEval(tt.input)
		testIntegerObject(t, ev, tt.expected)
	}
}

func testEval(input string) Object {
	p := NewParser(NewLexer(input))
	prg := p.ParseProgram()
	return Eval(prg, NewEnvironment())
}

func testIntegerObject(t *testing.T, o Object, expected int) bool {
	r, ok := o.(*Integer)
	if !ok {
		t.Errorf("obj is not integer. got=%T (%+v)", o, o)
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
		{"1< 2", true},
		{"1> 2", false},
		{"1< 1", false},
		{"1> 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 != 2", true},
		{"1 == 2", false},
	}
	for _, tt := range tests {
		t.Log("tt.input:", tt.input)
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

func TestIfElseEval(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (0) { 10 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, tt := range tests {
		evl := testEval(tt.input)
		i, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evl, i)
			continue
		}
		testNullObject(t, evl)
	}
}

func testNullObject(t *testing.T, o Object) bool {
	if o != NullObject {
		t.Errorf("object is not NULL. got=%T (%+v)", o, o)
		return false
	}
	return true
}

func TestReturnStmt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9", 10},
		{`
if (10 > 1) {
		if (10  > 1) {
			return 10;
		}
		return 1;
}`, 10},
	}
	for _, tt := range tests {
		evl := testEval(tt.input)
		testIntegerObject(t, evl, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct{ input, expected string }{
		{"5 + true", "unknown operator: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "unknown operator: INTEGER + BOOLEAN"},
		{"-true;", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{`
if (10 > 1) {
		if (10 > 1) {
			return true + false;
		}
		return 1;
}`, "unknown operator: BOOLEAN + BOOLEAN"},
	}
	for _, tt := range tests {
		evl := testEval(tt.input)
		testErrorCheck(t, evl, tt.expected)
	}
}

func testErrorCheck(t *testing.T, o Object, expected string) bool {
	err, ok := o.(*Error)
	if !ok {
		t.Errorf("not error object returned. got=%T (%+v)",
			o, o)
		return false
	}
	if err.Msg != expected {
		t.Errorf("wrong error message. expected=%q, got=%q",
			expected, err.Msg)
		return false
	}
	return true
}
