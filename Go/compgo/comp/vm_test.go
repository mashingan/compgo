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
