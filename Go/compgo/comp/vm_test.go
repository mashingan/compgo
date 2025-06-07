package comp

import (
	"compgo/interp"
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
		stackElm := vm.StackTop()
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
	}
}

func TestIntegerArithVm(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		// {"1 + 2", 3},
	}
	runVmTests(t, tests)
}
