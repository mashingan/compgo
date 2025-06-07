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
