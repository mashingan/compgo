package interp

import (
	"fmt"
	"testing"
)

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
		{"神業", "identifier not found: 神業"},
		{`"Hello" - "world"`, "unknown operator: STRING - STRING"},
		{`{"name": "Monkey"}[fn(x){ x }]`, "unknown as hash key: FUNCTION"},
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

func TestLetStatementEval(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let 主人公 = 5; 主人公;", 5},
		{"let 主人公 = 5 * 5; 主人公;", 25},
		{"let 主人公 = 5; let 友人 = 主人公; 友人;", 5},
		{"let 主人公 = 5; let 友人 = 主人公; let 神的 = 主人公 + 友人 + 5; 神的;", 15},
	}
	for _, tt := range tests {
		evl := testEval(tt.input)
		testIntegerObject(t, evl, tt.expected)
	}
}

func TestFunctionEval(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evl := testEval(input)
	fn, ok := evl.(*Function)
	if !ok {
		t.Fatalf("obj is not function. got=%T (%+v)", evl, evl)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameter. expected 1 got=%d", len(fn.Parameters))
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expBody := "(x+2)"
	if fn.Body.String() != expBody {
		t.Fatalf("body is not %q. got=%T", expBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let 自我 = fn(x) { x; }; 自我(5);", 5},
		{"let 自我 = fn(x) { return x; }; 自我(5);", 5},
		{"let 二倍 = fn(x) { x * 2; }; 二倍(5);", 10},
		{"let 誰何 = fn(役者, 女優) { 役者 + 女優; }; 誰何(5, 5);", 10},
		{`
let 誰何 = fn(役者, 女優) { 役者 + 女優; };
誰何(5+5, 誰何(5, 5));`, 20},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosure(t *testing.T) {
	input := `
let 新型 = fn(x) {
	fn(y) { x + y };
};

let 新型ｍｋ２ = 新型(2);
新型ｍｋ２(2);`
	testIntegerObject(t, testEval(input), 4)
}

func TestStringEval(t *testing.T) {
	exp := "hello 異世界!"
	input := fmt.Sprintf(`"%s"`, exp)
	evl := testEval(input)
	str, ok := evl.(*String)
	if !ok {
		t.Fatalf("obj is not String. got=%T (%+v)", evl, evl)
	}
	if str.Value != exp {
		t.Errorf("String has wrong value. got=%q, want=%q",
			str.Value, exp)
	}
}

func TestStringConcat(t *testing.T) {
	exp := "hello 異世界!"
	input := `"hello" + " " + "異世界!"`
	evl := testEval(input)
	str, ok := evl.(*String)
	if !ok {
		t.Fatalf("obj is not String. got=%T (%+v)", evl, evl)
	}
	if str.Value != exp {
		t.Errorf("String has wrong value. got=%q, want=%q",
			str.Value, exp)
	}
}

func TestBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello 異世界")`, 9},
		{`len(1)`, "argument to 'len' not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`first(["hello", "異世界"])`, "hello"},
		{`last(["hello", "異世界"])`, "異世界"},
		{`first("異世界")`, "異"},
		{`last("異世界")`, "界"},
		{`rest(["hello", "異", "世", "界"])`, `["異","世","界"]`},
		{`rest("hello 異世界")`, `ello 異世界`},
		{`push(["hello"], "異", "世", "界");`, `["hello", "異", "世", "界"]`},
		{`let a = ["hello"]; push(a, "異", "世", "界"); a;`,
			`["hello", "異", "世", "界"]`},
		{`push("hello ", "異", "世", "界");`, `hello 異世界`},
		{`let a = "hello "; push(a, "異", "世", "界"); a;`, `hello 異世界`},
	}
	for _, tt := range tests {
		evl := testEval(tt.input)
		switch exp := evl.(type) {
		case *Integer:
			texp, ok := tt.expected.(int)
			if !ok {
				t.Errorf("wrong expected type, want int got=%T (%+v)",
					tt.expected, tt.expected)
				continue
			}
			testIntegerObject(t, evl, texp)
		case *Error:
			sexp, _ := tt.expected.(string)
			if exp.Msg != sexp {
				t.Errorf("wrong error mesage. expected=%q, got=%q",
					exp, exp.Msg)
			}
		case *String:
			stxp, _ := tt.expected.(string)
			if exp.Value != stxp {
				t.Errorf("wrong string. expected=%q, got=%q",
					stxp, exp.Value)
			}

		}
	}
}

func TestSlice(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len([])`, 0},
		{`[1, 2, 3, 4 + 5]`, "[1,2,3,9]"},
		{`len([1, "two", true, "four"])`, 4},
		{`["hello", "異世界"]`, `["hello","異世界"]`},
		{`len(["世界一", "世界二"])`, 2},
	}
	for _, tt := range tests {
		evl := testEval(tt.input)
		switch exp := tt.expected.(type) {
		case int:
			testIntegerObject(t, evl, exp)
		case string:
			slc, ok := evl.(*SliceObj)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evl, evl)
				continue
			}
			if slc.Inspect() != exp {
				t.Errorf("wrong slice. expected=%q, got=%q",
					exp, slc.Inspect())
			}
		}
	}
}

func TestIndexEval(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`let a = [1, 2, 3, 4]; a[2]`, 3},
		{`[1, 2, 3, 4 + 5][3]`, 9},
		{`[][1]`, NullObject},
		{`["hello", "異世界"][2-1]`, "異世界"},
		{`let idx = 1*1; ["hello", "異世界"][idx]`, "異世界"},
		{`let idx = 3*2; "hello 異世界"[idx]`, "異"},
		{`let idx = 3*2; let isekai = "hello 異世界"; isekai[idx]`, "異"},
	}
	for _, tt := range tests {
		evl := testEval(tt.input)
		switch exp := tt.expected.(type) {
		case int:
			testIntegerObject(t, evl, exp)
		case string:
			str, ok := evl.(*String)
			if !ok {
				t.Errorf("object is not string. got=%T (%+v)", evl, evl)
				continue
			}
			if str.Value != exp {
				t.Errorf("wrong string. expected=%q, got=%q",
					exp, str.Value)
			}
		case *Null:
			testNullObject(t, evl)
		}
	}
}

func TestHashLiteralEval(t *testing.T) {
	input := `let two = "two";
{
	"one": 10 - 9,
	two: 1 + 1,
	"thr" + "ee" : 6 / 2,
	4: 4,
	true: 5,
	false: 6,
}`
	evl := testEval(input)
	r, ok := evl.(*Hash)
	if !ok {
		t.Fatalf("object is not hash. got=%T (%+v)", evl, evl)
	}
	exps := map[HashKey]int{
		(&String{Primitive[string]{"one"}}).HashKey():   1,
		(&String{Primitive[string]{"two"}}).HashKey():   2,
		(&String{Primitive[string]{"three"}}).HashKey(): 3,
		(&Integer{Primitive[int]{4}}).HashKey():         4,
		TrueObject.HashKey():                            5,
		FalseObject.HashKey():                           6,
	}
	if len(r.Pairs) != len(exps) {
		t.Fatalf("got wrong number of pairs. got=%d", len(r.Pairs))
	}

	for expk, expv := range exps {
		pair, ok := r.Pairs[expk]
		if !ok {
			t.Error("no pair given key in pairs")
		}
		testIntegerObject(t, pair.Value, expv)
	}

}

func TestHashIndexEval(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`let key = "foo"; {"foo": 5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5:5}[5]`, 5},
		{`{true:5}[true]`, 5},
		{`{false:5}[false]`, 5},
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

func TestQuoteEval(t *testing.T) {
	tests := []struct{ input, expected string }{
		{"quote(5)", "5"},
		{"quote(5+8)", "(5+8)"},
		{"quote(異世界)", "異世界"},
		{"quote(異世界+勇者)", "(異世界+勇者)"},
		{"let 勇者 = 8; quote(勇者)", "勇者"},
	}
	for _, tt := range tests {
		evl := testEval(tt.input)
		quote, ok := evl.(*Quote)
		if !ok {
			t.Fatalf("expected quote. got=%T (%+v)", evl, evl)
		}
		if quote.Node == nil {
			t.Fatal("quote.Node is nil")
		}
		if quote.Node.String() != tt.expected {
			t.Errorf("not equal. got=%q want=%q", quote.Node.String(), tt.expected)
		}
	}
}

func TestUnquoteEval(t *testing.T) {
	tests := []struct{ input, expected string }{
		{"quote(unquote(4))", "4"},
		{"quote(unquote(4+4))", "8"},
		{"quote(8 + unquote(4+4))", "(8+8)"},
		{"quote(unquote(4+4) + 8)", "(8+8)"},
		{"let 勇者 = 8; quote(unquote(勇者))", "8"},
		{"quote(unquote(true))", "true"},
		{"quote(unquote(true == false))", "false"},
		{"quote(unquote(quote(4 + 4)))", "(4+4)"},
		{"let quotedInfix = quote(4 + 4); quote(unquote(4 + 4) + unquote(quotedInfix))",
			"(8+(4+4))"},
	}
	for _, tt := range tests {
		evl := testEval(tt.input)
		q, ok := evl.(*Quote)
		if !ok {
			t.Errorf("expected *quote. got=%T (%+v)", evl, evl)
			continue
		}
		if q.Node == nil {
			t.Error("quote.Node is nil")
		}
		if q.Node.String() != tt.expected {
			t.Errorf("not equal. got=%q want=%q", q.Node.String(), tt.expected)
		}
	}
}

func TestDefineMacro(t *testing.T) {
	input := `
let number = 1;
let func = fn(x, y) { x + y; };
let mymacro = macro(x, y) { x + y; };
`
	p := NewParser(NewLexer(input))
	prg := p.ParseProgram()
	env := NewEnvironment()
	DefineMacros(prg, env)
	if len(prg.Statements) != 2 {
		t.Fatalf("wrong number of statements. got=%d", len(prg.Statements))
	}
	_, ok := env.Get("number")
	if ok {
		t.Fatal("number should not defined")
	}
	_, ok = env.Get("func")
	if ok {
		t.Fatal("func should not defined")
	}
	m, ok := env.Get("mymacro")
	if !ok {
		t.Fatal("mymacro not in environment")
	}
	macro, ok := m.(*MacroObj)
	if !ok {
		t.Fatalf("object is not macro. got=%T (%+v)", m, m)
	}
	if len(macro.Parameters) != 2 {
		t.Fatalf("wrong number of statements. got=%d", len(prg.Statements))
	}
	testLiteralExpression(t, macro.Parameters[0], "x")
	testLiteralExpression(t, macro.Parameters[1], "y")
	expbody := "(x+y)"
	if macro.Body.String() != expbody {
		t.Errorf("macro body is not expected=%q, got=%q",
			expbody, macro.Body)
	}
}

func TestExpandMacro(t *testing.T) {
	tests := []struct{ input, expected string }{
		{
			`let infix = macro() { quote(1+2); };infix();`,
			`(1+2)`,
		},
		{
			`let rev = macro(a, b) { quote(unquote(b) - unquote(a)); };
			rev(2+2, 10-5);`,
			`(10-5)-(2+2)`,
		},
	}
	for _, tt := range tests {
		p := NewParser(NewLexer(tt.expected))
		prgexp := p.ParseProgram()
		p = NewParser(NewLexer(tt.input))
		prginp := p.ParseProgram()
		env := NewEnvironment()
		DefineMacros(prginp, env)
		t.Log("prginp:", prginp)
		expanded := ExpandMacros(prginp, env)
		if expanded.String() != prgexp.String() {
			t.Errorf("not equal. want=%q got=%q",
				prgexp.String(), expanded.String())
		}
	}
}
