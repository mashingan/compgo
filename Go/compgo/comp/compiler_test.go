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
			Make(OpAdd),
			Make(OpPop),
		}},
		{"1;2", []any{1, 2}, []Instructions{
			Make(OpConstant, 0),
			Make(OpPop),
			Make(OpConstant, 1),
			Make(OpPop),
		}},
		{"1-2", []any{1, 2}, []Instructions{
			Make(OpConstant, 0),
			Make(OpConstant, 1),
			Make(OpSub),
			Make(OpPop),
		}},
		{"1*2", []any{1, 2}, []Instructions{
			Make(OpConstant, 0),
			Make(OpConstant, 1),
			Make(OpMul),
			Make(OpPop),
		}},
		{"1/2", []any{1, 2}, []Instructions{
			Make(OpConstant, 0),
			Make(OpConstant, 1),
			Make(OpDiv),
			Make(OpPop),
		}},
		{"-1", []any{1}, []Instructions{
			Make(OpConstant, 0),
			Make(OpMinus),
			Make(OpPop),
		}},
	}
	runCompilerTest(t, tests)
}

func TestBooleanCompile(t *testing.T) {
	tests := []compilerTestCase{
		{"1==2", []any{1, 2}, []Instructions{
			Make(OpConstant, 0),
			Make(OpConstant, 1),
			Make(OpEq),
			Make(OpPop),
		}},
		{"1!=2", []any{1, 2}, []Instructions{
			Make(OpConstant, 0),
			Make(OpConstant, 1),
			Make(OpNeq),
			Make(OpPop),
		}},
		{"1>2", []any{1, 2}, []Instructions{
			Make(OpConstant, 0),
			Make(OpConstant, 1),
			Make(OpGt),
			Make(OpPop),
		}},
		{"1<2", []any{1, 2}, []Instructions{
			Make(OpConstant, 0),
			Make(OpConstant, 1),
			Make(OpLt),
			Make(OpPop),
		}},
		{"1>=2", []any{1, 2}, []Instructions{
			Make(OpConstant, 0),
			Make(OpConstant, 1),
			Make(OpGte),
			Make(OpPop),
		}},
		{"1<=2", []any{1, 2}, []Instructions{
			Make(OpConstant, 0),
			Make(OpConstant, 1),
			Make(OpLte),
			Make(OpPop),
		}},
		{"true", []any{}, []Instructions{
			Make(OpTrue),
			Make(OpPop),
		}},
		{"false", []any{}, []Instructions{
			Make(OpFalse),
			Make(OpPop),
		}},
		{"!false", []any{}, []Instructions{
			Make(OpFalse),
			Make(OpBang),
			Make(OpPop),
		}},
		{"!true", []any{}, []Instructions{
			Make(OpTrue),
			Make(OpBang),
			Make(OpPop),
		}},
	}
	runCompilerTest(t, tests)
}

func parse(input string) *interp.Program {
	p := interp.NewParser(interp.NewLexer(input))
	return p.ParseProgram()
}

func testInstructions(t *testing.T, expected []Instructions, got Instructions) error {
	exp := Instructions{}
	for _, ex := range expected {
		exp = append(exp, ex...)
	}
	if len(exp) != len(got) {
		return fmt.Errorf("wrong instruction length.\ngot=%q\nwant=%q", got, exp)
	}
	for i, b := range exp {
		if b != got[i] {
			t.Logf("expected: %q\n", exp)
			t.Logf("got: %q", got)
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
		if err := testConstants(t, tt.expectedConstants, bc.Constants); err != nil {
			t.Error(err)
			continue
		}
		err = testInstructions(t, tt.expectedInstructions, bc.Instructions)
		if err != nil {
			t.Fatalf("test instruction failed: %s", err)
		}

	}

}

func testConstants(t *testing.T, expected []any, actual []interp.Object) error {
	if len(expected) != len(actual) {
		t.Logf("got=%q\nwant=%q", actual, expected)
		return fmt.Errorf("wrong number of constants. got=%d want=%d", len(actual), len(expected))
	}
	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		case []Instructions:
			fn, ok := actual[i].(*CompiledFunction)
			if !ok {
				return fmt.Errorf("constant %d - not a function: %T", i, actual[i])
			}
			if err := testInstructions(t, constant, fn.Instructions); err != nil {
				return fmt.Errorf("constant %d - testInstruction failed: %s", i, err)

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
		return fmt.Errorf("object wrong value. got=%d want=%d", i.Value, n)
	}

	return nil
}

func TestInstructionsString(t *testing.T) {
	inst := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpPop),
	}
	t.Log("1:", Make(OpConstant, 1))
	t.Log("2:", Make(OpConstant, 2))
	t.Log("3:", Make(OpConstant, 65535))
	expected := strings.TrimSpace(`
0000 OpAdd
0001 OpConstant 1
0004 OpConstant 2
0007 OpConstant 65535
0010 OpPop
`)
	insts := Instructions{}
	for _, ins := range inst {
		insts = append(insts, ins...)
	}
	if insts.String() != expected {
		t.Errorf("instruction wrong format.\nwant=%q\ngot=%q", expected, insts.String())
	}
}

func TestConditionalsCompile(t *testing.T) {
	tests := []compilerTestCase{
		{`if (true) { 10 }; 3333`, []any{10, 3333}, []Instructions{
			Make(OpTrue),            // 0000
			Make(OpJumpIfFalsy, 10), // 0001
			Make(OpConstant, 0),     // 0004
			Make(OpJump, 11),        // 0007
			Make(OpNull),            // 0010
			Make(OpPop),             // 0011
			Make(OpConstant, 1),     // 0012
			Make(OpPop),             // 0015
		}},
		{`if (true) { 10 } else { 20 }; 3333`, []any{10, 20, 3333}, []Instructions{
			Make(OpTrue),            // 0000
			Make(OpJumpIfFalsy, 10), // 0001
			Make(OpConstant, 0),     // 0004
			Make(OpJump, 13),        // 0007
			Make(OpConstant, 1),     // 0010
			Make(OpPop),             // 00013
			Make(OpConstant, 2),     // 00014
			Make(OpPop),             // 00017
		}},
		{`if (false) { 10 }`, []any{10}, []Instructions{
			Make(OpFalse),           // 0000
			Make(OpJumpIfFalsy, 10), // 0001
			Make(OpConstant, 0),     // 0004
			Make(OpJump, 11),        // 0007
			Make(OpNull),            // 0010
			Make(OpPop),             // 0011
		}},
	}
	runCompilerTest(t, tests)
}
func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
		let one = 1;
		let two = 2;
		`,
			expectedConstants: []any{1, 2},
			expectedInstructions: []Instructions{
				Make(OpConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpConstant, 1),
				Make(OpSetGlobal, 1),
			},
		},
		{
			input: `
		let one = 1;
		one;
		`,
			expectedConstants: []any{1},
			expectedInstructions: []Instructions{
				Make(OpConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpGetGlobal, 0),
				Make(OpPop),
			},
		},
		{
			input: `
		let one = 1;
		let two = one;
		two;
		`,
			expectedConstants: []any{1},
			expectedInstructions: []Instructions{
				Make(OpConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpGetGlobal, 0),
				Make(OpSetGlobal, 1),
				Make(OpGetGlobal, 1),
				Make(OpPop),
			},
		},
	}
	runCompilerTest(t, tests)
}

func TestStringExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"monkey"`,
			expectedConstants: []any{"monkey"},
			expectedInstructions: []Instructions{
				Make(OpConstant, 0),
				Make(OpPop),
			},
		},
		{
			input:             `"異" + "世界"`,
			expectedConstants: []any{"異", "世界"},
			expectedInstructions: []Instructions{
				Make(OpConstant, 0),
				Make(OpConstant, 1),
				Make(OpAdd),
				Make(OpPop),
			},
		},
	}
	runCompilerTest(t, tests)
}

func TestArrayExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `[]`,
			expectedConstants: []any{},
			expectedInstructions: []Instructions{
				Make(OpArray, 0),
				Make(OpPop),
			},
		},
		{
			input:             `["異", "世界"]`,
			expectedConstants: []any{"異", "世界"},
			expectedInstructions: []Instructions{
				Make(OpConstant, 0),
				Make(OpConstant, 1),
				Make(OpArray, 2),
				Make(OpPop),
			},
		},
		{
			input:             `[1 + 2, 3 - 4, 5 * 6]`,
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []Instructions{
				Make(OpConstant, 0),
				Make(OpConstant, 1),
				Make(OpAdd),
				Make(OpConstant, 2),
				Make(OpConstant, 3),
				Make(OpSub),
				Make(OpConstant, 4),
				Make(OpConstant, 5),
				Make(OpMul),
				Make(OpArray, 3),
				Make(OpPop),
			},
		},
	}
	runCompilerTest(t, tests)
}

func TestHashExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `{}`,
			expectedConstants: []any{},
			expectedInstructions: []Instructions{
				Make(OpHash, 0),
				Make(OpPop),
			},
		},
		{
			input:             `{"i": "異", "sekai": "世界"}`,
			expectedConstants: []any{"i", "異", "sekai", "世界"},
			expectedInstructions: []Instructions{
				Make(OpConstant, 0),
				Make(OpConstant, 1),
				Make(OpConstant, 2),
				Make(OpConstant, 3),
				Make(OpHash, 4),
				Make(OpPop),
			},
		},
		{
			input:             `{1: 2, 3: 4, 5: 6}`,
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []Instructions{
				Make(OpConstant, 0),
				Make(OpConstant, 1),
				Make(OpConstant, 2),
				Make(OpConstant, 3),
				Make(OpConstant, 4),
				Make(OpConstant, 5),
				Make(OpHash, 6),
				Make(OpPop),
			},
		},
		{
			input:             `{1: 2 + 3, 4: 5 * 6}`,
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []Instructions{
				Make(OpConstant, 0),
				Make(OpConstant, 1),
				Make(OpConstant, 2),
				Make(OpAdd),
				Make(OpConstant, 3),
				Make(OpConstant, 4),
				Make(OpConstant, 5),
				Make(OpMul),
				Make(OpHash, 4),
				Make(OpPop),
			},
		},
	}
	runCompilerTest(t, tests)
}

func TestIndexExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `[1, 2, 3][1+1]`,
			expectedConstants: []any{1, 2, 3, 1, 1},
			expectedInstructions: []Instructions{
				Make(OpConstant, 0),
				Make(OpConstant, 1),
				Make(OpConstant, 2),
				Make(OpArray, 3),
				Make(OpConstant, 3),
				Make(OpConstant, 4),
				Make(OpAdd),
				Make(OpIndex),
				Make(OpPop),
			},
		},
		{
			input:             `{1: 2}[2-1]`,
			expectedConstants: []any{1, 2, 2, 1},
			expectedInstructions: []Instructions{
				Make(OpConstant, 0),
				Make(OpConstant, 1),
				Make(OpHash, 2),
				Make(OpConstant, 2),
				Make(OpConstant, 3),
				Make(OpSub),
				Make(OpIndex),
				Make(OpPop),
			},
		},
	}
	runCompilerTest(t, tests)
}
