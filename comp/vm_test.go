package comp

import (
	"compgo/interp"
	"fmt"
	"testing"
)

type vmTestCase struct {
	input    string
	expected any
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()
	for _, tt := range tests {
		prg := parse(tt.input)
		comp := New()
		err := comp.Compile(prg)
		if err != nil {
			t.Fatalf("compile error: %s", err)
		}
		vm := NewVm(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}
		stackElm := vm.LastPop()
		testExpectedObject(t, tt.expected, stackElm)
	}
}

func testExpectedObject(t *testing.T, expected any, actual interp.Object) {
	t.Helper()
	switch exp := expected.(type) {
	case int:
		err := testIntegerObject(exp, actual)
		if err != nil {
			t.Errorf("integer object test fail: %s", err)
		}
	case bool:
		err := testBooleanObject(exp, actual)
		if err != nil {
			t.Errorf("%s", err)
		}
	case nil:
		if actual != interp.NullObject {
			t.Errorf("object is not null. got=%T (%+v)", actual, actual)
		}
	case string:
		err := testStringObject(exp, actual)
		if err != nil {
			t.Errorf("%s", err)
		}
	case []any:
		if err := testArrayObject(exp, actual); err != nil {
			t.Errorf("%s", err)
		}
	case map[interp.HashKey]int:
		h, ok := actual.(*interp.Hash)
		if !ok {
			t.Errorf("object is not hash. got=%T (%+v)", actual, actual)
		}
		for k, v := range exp {
			kv, ok := h.Pairs[k]
			if !ok {
				t.Errorf("key hash %v not in map", k)
				continue
			}
			switch kvk := kv.Key.(type) {
			case *interp.Integer:
				if kvk.HashKey() != k {
					t.Errorf("key %d is not same hash.", kvk.Value)
					continue
				}
			}
			switch kvv := kv.Value.(type) {
			case *interp.Integer:
				if kvv.Value != v {
					t.Errorf("integer value is not same. got=%d want=%d",
						kvv.Value, v)
				}
			}
		}
	case *interp.Error:
		err, ok := actual.(*interp.Error)
		if !ok {
			t.Errorf("object is not error. got=%T (%+v)", actual, actual)
		} else if err.Msg != exp.Msg {
			t.Errorf("error message is wrong. got=%q want=%q",
				err.Msg, exp.Msg)

		}
	}
}

func testBooleanObject(expected bool, actual interp.Object) error {
	b, ok := actual.(*interp.Boolean)
	if !ok {
		return fmt.Errorf("object is not boolean. got=%T (%+v)", actual, actual)
	}
	if b.Value != expected {
		return fmt.Errorf("object is wrong value. got=%t want=%t", b.Value, expected)
	}
	return nil
}

func testStringObject(expected string, actual interp.Object) error {
	s, ok := actual.(*interp.String)
	if !ok {
		return fmt.Errorf("object is not string. got=%T (%+v)", actual, actual)
	}
	if s.Value != expected {
		return fmt.Errorf("object is wrong value. got=%s want=%s", s.Value, expected)
	}
	return nil
}

func testArrayObject(expected []any, actual interp.Object) error {
	arr, ok := actual.(*interp.SliceObj)
	if !ok {
		return fmt.Errorf("object is not array. got=%T (%+v)", actual, actual)
	}
	if len(expected) != len(arr.Elements) {
		return fmt.Errorf("array length is not match. got=%d want=%d",
			len(arr.Elements), len(expected))
	}
	for i, e := range expected {
		switch e := e.(type) {
		case bool:
			return testBooleanObject(e, arr.Elements[i])
		case int:
			return testIntegerObject(e, arr.Elements[i])
		case string:
			return testStringObject(e, arr.Elements[i])
		case nil:
			if arr.Elements[i] != interp.NullObject {
				return fmt.Errorf("object is not null.")
			}
		}
	}
	return nil
}

func TestIntegerArithVm(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 * 2", 2},
		{"2 / 2", 1},
		{"3 - 2", 1},
		{"-5", -5},
		{"-1 * 10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	runVmTests(t, tests)
}

func TestBooleanVm(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"true == true", true},
		{"false == false", true},
		{"2 == 2", true},
		{"2 != 2", false},
		{"2 > 3", false},
		{"2 < 3", true},
		{"2 <= 3", true},
		{"2 >= 3", false},
		{"(1 < 2) == true", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	runVmTests(t, tests)
}

func TestConditionalsVm(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 } else { 20 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (false) { 10 }", nil},
		{"if (1 > 2) { 10 }", nil},
	}
	runVmTests(t, tests)
}

func TestGlobalLetStatementsVm(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one;", 1},
		{"let one = 1; one + one;", 2},
		{"let one = 1; let two = 2; one + two;", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}
	runVmTests(t, tests)
}

func TestStringVm(t *testing.T) {
	tests := []vmTestCase{
		{`let monkey = "monkey"; monkey;`, "monkey"},
		{`let i = "異"; let sekai = "世界"; i + sekai;`, "異世界"},
		{`let i = "異"; let sekai = "世界"; let isekai = i + sekai; isekai`, "異世界"},
	}
	runVmTests(t, tests)
}

func TestArrayVm(t *testing.T) {
	tests := []vmTestCase{
		{`[]`, []int{}},
		{`["異世界", 1, "勇者", 2]`, []any{"異世界", 1, "勇者", 2}},
		{`let i = "異"; let kai = "界"; [i + "世" + kai, 10 / 10, "勇" + "者", 10 * 10 / 50]`,
			[]any{"異世界", 1, "勇者", 2}},
		{`["異世界", 1] + ["勇者", 2]`, []any{"異世界", 1, "勇者", 2}},
		{`["異世界", 1, "hehe"] + ["勇者", 2, "hallo"]`, []any{"異世界", 1, "hehe", "勇者", 2, "hallo"}},
	}
	runVmTests(t, tests)
}

func TestHashVm(t *testing.T) {
	tests := []vmTestCase{
		{`{}`, map[interp.HashKey]any{}},
		{`{1: 2, 2: 3}`, map[interp.HashKey]int{
			(&interp.Integer{Primitive: interp.Primitive[int]{Value: 1}}).HashKey(): 2,
			(&interp.Integer{Primitive: interp.Primitive[int]{Value: 2}}).HashKey(): 3,
		}},
		{`{1+1: 2*2, 3+3: 4*4}`, map[interp.HashKey]int{
			(&interp.Integer{Primitive: interp.Primitive[int]{Value: 2}}).HashKey(): 4,
			(&interp.Integer{Primitive: interp.Primitive[int]{Value: 6}}).HashKey(): 16,
		}},
	}
	runVmTests(t, tests)
}

func TestIndexVm(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0+2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", nil},
		{"[1, 2, 3][99]", nil},
		{"[1][-1]", nil},
		{"{1:1, 2:2}[1]", 1},
		{"{1:1, 2:2}[2]", 2},
		{"{}[2]", nil},
		{"{1:1, 2:2}[0]", nil},
		{`""[-1]`, nil},
		{`"異世界"[1]`, "世"},
		{`"異世界"[-1]`, nil},
	}
	runVmTests(t, tests)
}

func TestFunctionVm_call(t *testing.T) {
	tests := []vmTestCase{
		{"let fpt = fn() { 5 + 10 }; fpt();", 15},
		{`
		let one = fn() { 1 }; 
		let two = fn() { 2 };
		one() + two();`, 3},
		{`
		let a = fn() { 1 }; 
		let b = fn() { a() + 1 };
		let c = fn() { b() + 1 };
		c();`, 3},
		{`
		let early = fn() { return 99; 100; }; 
		early();`, 99},
		{`
		let early = fn() { return 99; return 100; }; 
		early();`, 99},
		{`
		let noreturn = fn() {}; 
		noreturn();`, nil},
		{`
		let noreturn = fn() {}; 
		let noret = fn() { noreturn() };
		noret();`, nil},
	}
	runVmTests(t, tests)
}

func TestFunctionVm_firstClass(t *testing.T) {
	tests := []vmTestCase{
		{`
		let retone = fn() { 1 }; 
		let ret2one = fn() { retone; };
		ret2one()()`, 1},
		{`
		let retone = fn() {
			let r1 = fn() { 1 };
			r1
		}; 
		retone()()`, 1},
	}
	runVmTests(t, tests)
}

func TestFunctionVm_localBinding(t *testing.T) {
	tests := []vmTestCase{
		{`let one = fn() { let one = 1; one };
		one();`, 1},
		{`
		let one2 = fn() {
			let one = 1;
			let two = 2;
			one + two;
		};
		one2();`, 3},
		{`
		let one2 = fn() {
			let one = 1;
			let two = 2;
			one + two;
		};
		let three4 = fn() {
			let three = 3;
			let four = 4;
			three + four;
		};
		one2() + three4();`, 10},
		{`
		let f1 = fn() {
			let foobar = 50;
			foobar;
		};
		let f2 = fn() {
			let foobar = 100;
			foobar;
		};
		f1() + f2();`, 150},
		{`
		let globalseed = 50;
		let m1 = fn() {
			let num = 1;
			globalseed - num;
		};
		let m2 = fn() {
			let num = 2;
			globalseed - num;
		};
		m1() + m2();`, 97},
	}
	runVmTests(t, tests)
}

func TestFunctionVm_argsBinding(t *testing.T) {
	tests := []vmTestCase{
		{`let identity = fn(a) { a; }; identity(4);`, 4},
		{`let sum = fn(a, b) { let c = a + b; c; }; sum(1, 2);`, 3},
		{`let sum = fn(a, b) { let c = a - b; c; }; sum(1, 2);`, -1},
		{`let sum = fn(a, b) { let c = a + b; c; }; sum(1, 2) + sum(3, 4);`, 10},
		{`let glob = 10;
		let sum = fn(a, b) {
			let c = a + b;
			c + glob;
		};
		let outer = fn() {
			sum(1, 2) + sum(3, 4) + glob;
		};
		outer() + glob;`, 50},
	}
	runVmTests(t, tests)
}

func TestFunctionVm_wrongArgNum(t *testing.T) {
	tests := []vmTestCase{
		{`fn(){1;}(1)`, `wrong argument number: want=0, got=1`},
		{`fn(a){a;}()`, `wrong argument number: want=1, got=0`},
		{`fn(a, b){ a + b;}(1)`, `wrong argument number: want=2, got=1`},
	}
	for _, tt := range tests {
		prg := parse(tt.input)
		comp := New()
		_ = comp.Compile(prg)
		vm := NewVm(comp.Bytecode())
		err := vm.Run()
		if err == nil {
			t.Fatal("expected vm but resulted none")
		}
		if err.Error() != tt.expected {
			t.Fatalf("wrong vm error message: want=%q got=%q", tt.expected, err.Error())
		}
	}
}

func TestBuiltinVm(t *testing.T) {
	tests := []vmTestCase{
		{`len("")`, 0},
		{`len([])`, 0},
		{`len("異世界")`, 3},
		{`len(["異", "世", "界"])`, 3},
		{`len(1)`, &interp.Error{
			Msg: "argument to 'len' not supported, got INTEGER"}},
		{`len("異世界", "勇者")`, &interp.Error{
			Msg: "wrong number of arguments. got=2, want=1",
		}},
		{`puts("異世界", "勇者")`, nil},
		{`first("異世界")`, "異"},
		{`first(["異", "世", "界"])`, "異"},
		{`first(1)`, &interp.Error{
			Msg: "argument to 'first' not supported, got INTEGER"}},
		{`first("異世界", "勇者")`, &interp.Error{
			Msg: "wrong number of arguments. got=2, want=1",
		}},
		{`first("")`, nil},
		{`first([])`, nil},
		{`last("異世界")`, "界"},
		{`last(["異", "世", "界"])`, "界"},
		{`last(1)`, &interp.Error{
			Msg: "argument to 'last' not supported, got INTEGER"}},
		{`last("異世界", "勇者")`, &interp.Error{
			Msg: "wrong number of arguments. got=2, want=1",
		}},
		{`last("")`, nil},
		{`last([])`, nil},
		{`rest("異世界")`, "世界"},
		{`rest(["異", "世", "界"])`, []any{"世", "界"}},
		{`rest(1)`, &interp.Error{
			Msg: "argument to 'rest' not supported, got INTEGER"}},
		{`rest("異世界", "勇者")`, &interp.Error{
			Msg: "wrong number of arguments. got=2, want=1",
		}},
		{`rest("")`, nil},
		{`rest([])`, nil},
		{`push("異世界")`, "異世界"},
		{`push("異世", "界")`, "異世界"},
		{`push(["異", "世", "界"])`, []any{"異", "世", "界"}},
		{`push(["異", "世"], "界")`, []any{"異", "世", "界"}},
		{`push(1, 1)`, &interp.Error{
			Msg: "argument to 'push' not supported, got INTEGER"}},
		{`push()`, &interp.Error{
			Msg: "wrong number of arguments. got=0, want=2",
		}},
	}
	runVmTests(t, tests)
}

func TestClosureVm(t *testing.T) {
	tests := []vmTestCase{
		{`let adder = fn(base) {
			let inner = fn(add) {
				base + add;
			};
			inner
		}; adder(10)(5)`, 15},
		{`let newc = fn(a) {
			fn() { a; };
		};
		let closure = newc(99);
		closure()`, 99},
	}
	runVmTests(t, tests)
}

func TestClosureVm_recursive(t *testing.T) {
	tests := []vmTestCase{
		{`let countdown = fn(x) {
			if (x == 0) {
				return 0;
			} else {
				countdown(x - 1);
			}
		}; countdown(2)`, 0},
		{`let countdown = fn(x) {
			if (x == 0) {
				return 0;
			} else {
				countdown(x - 1);
			}
		};
		let wrapper = fn() {
			countdown(2);
		};
		wrapper();`, 0},
		{`let wrapper = fn() {
			let countdown = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					countdown(x - 1);
				}
			};
			countdown(1);
		};
		wrapper();`, 0},
	}
	runVmTests(t, tests)
}

func TestRecusiveFibonacci(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
let fibonacci = fn(x) {
	if (x == 0) { return 0 }
	else {
		if (x == 1) { return 1; }
		else { fibonacci(x - 1) + fibonacci(x - 2) }
	}
};
fibonacci(15);
			`,
			expected: 610,
		},
	}
	runVmTests(t, tests)
}
