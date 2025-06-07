package comp

import (
	"compgo/interp"
	"fmt"
	"strings"
	"testing"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []Instructions
}

func TestIntegerArith(t *testing.T) {
	tests := []compilerTestCase{
		{"1+2", []any{1, 2}, []Instructions{
			Make(OpConstant, 0),
			Make(OpConstant, 1),
		}},
	}
	runCompilerTest(t, tests)
}

func parse(input string) *interp.Program {
	p := interp.NewParser(interp.NewLexer(input))
	return p.ParseProgram()
}

func testInstructions(expected []Instructions, got Instructions) error {
	exp := Instructions{}
	for _, ex := range expected {
		exp = append(exp, ex...)
	}
	if len(exp) != len(got) {
		return fmt.Errorf("wrong instruction length.\ngot=%q\nwant=%q", got, exp)
	}
	for i, b := range exp {
		if b != got[i] {
			return fmt.Errorf("wrong byte at instruction byte %d. got=%d want=%d",
				i, b, got[i])
		}
	}
	return nil
}

func runCompilerTest(t *testing.T, ct []compilerTestCase) {
	t.Helper()
	for _, tt := range ct {
		prg := parse(tt.input)
		compiler := New()
		err := compiler.Compile(prg)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}
		bc := compiler.Bytecode()
		err = testInstructions(tt.expectedInstructions, bc.Instructions)
		if err != nil {
			t.Fatalf("test instruction failed: %s", err)
		}

	}

}

func testConstants(t testing.T, expected []any, actual []interp.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d want=%d", len(actual), len(expected))
	}
	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}
	return nil
}

func testIntegerObject(n int, o interp.Object) error {
	i, ok := o.(*interp.Integer)
	if !ok {
		return fmt.Errorf("object is not integer. got=%T (%+v)", o, o)
	}
	if i.Value != n {
		return fmt.Errorf("object wrong integer. got=%d want=%d", i.Value, n)
	}

	return nil
}

func TestInstructionsString(t *testing.T) {
	inst := []Instructions{
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}
	t.Log("1:", Make(OpConstant, 1))
	t.Log("2:", Make(OpConstant, 2))
	t.Log("3:", Make(OpConstant, 65535))
	expected := strings.TrimSpace(`
0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
`)
	insts := Instructions{}
	for _, ins := range inst {
		insts = append(insts, ins...)
	}
	if insts.String() != expected {
		t.Errorf("instruction wrong format.\nwant=%q\ngot=%q", expected, insts.String())
	}
}
