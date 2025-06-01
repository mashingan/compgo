package interp

import (
	"fmt"
	"testing"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input, expectedIdentifier string
		expectedValue             any
	}{
		{"let x一 = 5;", "x一", 5},
		{"let y二 = 10;", "y二", 10},
		{"let foobar=838383;", "foobar", 838383},
		{"let 剣士 =　前衛;", "剣士", "前衛"},
		{"let 魔法使い = true;", "魔法使い", true},
	}
	for _, tt := range tests {
		l := NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got %d",
				len(program.Statements))
		}
		stmt := program.Statements[0]
		testLetStatement(t, stmt, tt.expectedIdentifier)
		val := stmt.(*LetStatement).Value
		testLiteralExpression(t, val, tt.expectedValue)
	}
}

func testLetStatement(t *testing.T, s Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf(`TokenLiteral not "let". got=%q\n`, s.TokenLiteral())
		return false
	}
	letstmt, ok := s.(*LetStatement)
	if !ok {
		t.Errorf("not statement. got=%T\n", s)
		return false
	}
	if letstmt.Name.Value != name {
		t.Errorf(`let stmt value expected "%s". got="%s"`, name,
			letstmt.Name.Value)
		return false
	}
	if letstmt.Name.TokenLiteral() != name {
		t.Errorf(`let stmt name expected "%s". got="%s"`, name,
			letstmt.Name)
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errs))
	for _, msg := range errs {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatment(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	p := NewParser(NewLexer(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		retstmt, ok := stmt.(*ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ReturnStatement. got=%T", stmt)
			continue
		}
		if retstmt.TokenLiteral() != "return" {
			t.Errorf("stmt token literal not 'return', got=%q,",
				retstmt.TokenLiteral())
		}
	}

}

func testIdentifier(t *testing.T, exp Expression, value string) bool {
	ident, ok := exp.(*Identifier)
	if !ok {
		t.Errorf("exp (%s) is not identifier. got=%T", exp, exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("identifier expected '%s'. got='%s'", value, ident.Value)
		return false
	}
	if ident.Literal != value {
		t.Errorf("identifier literal expected '%s'. got='%s'", value, ident.Value)
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntLiteral(t, exp, v)
	case int8, int32, int64, uint, uint8, uint16, uint32, uint64:
		vv := v.(int)
		return testIntLiteral(t, exp, vv)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanExpression(t, exp, v)
	}
	return false
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar異世界"
	p := NewParser(NewLexer(input))
	prog := p.ParseProgram()
	checkParserErrors(t, p)
	if len(prog.Statements) != 1 {
		t.Fatalf("wrong program statements. got=%d", len(prog.Statements))
	}
	stmt, ok := prog.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not expression statement. got=%T", prog.Statements[0])
	}
	testIdentifier(t, stmt.Expression, input)
}

func TestIntLiteralExpression(t *testing.T) {
	input := "5;"
	expected := 5
	p := NewParser(NewLexer(input))
	prog := p.ParseProgram()
	checkParserErrors(t, p)
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 program statments. got=%d", len(prog.Statements))
	}
	stmt, ok := prog.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not expression statement. got=%T", prog.Statements[0])
	}
	testIntLiteral(t, stmt.Expression, expected)
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		val      int
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}
	for _, tt := range prefixTests {
		p := NewParser(NewLexer(tt.input))
		prog := p.ParseProgram()
		checkParserErrors(t, p)
		if len(prog.Statements) != 1 {
			t.Fatalf("expected 1 program statments. got=%d", len(prog.Statements))
		}
		stmt, ok := prog.Statements[0].(*ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not expression statement. got=%T", prog.Statements[0])
		}
		exp, ok := stmt.Expression.(*PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not prefix expression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.operator is not '%s'. got='%s'",
				tt.operator, exp.Operator)
		}
		testIntLiteral(t, exp.Right, tt.val)
	}
}

func testIntLiteral(t *testing.T, ex Expression, val int) bool {
	itg, ok := ex.(*IntLiteral)
	if !ok {
		t.Errorf("ex is not *IntLiteral. got=%T", ex)
		return false
	}
	if itg.Value != val {
		t.Errorf("itg.value not %d. got=%d", val, itg.Value)
		return false
	}
	if itg.Literal != fmt.Sprintf("%d", val) {
		t.Errorf("itg.literal not %d. got=%s", val, itg.Literal)
		return false
	}
	return true
}

func TestParsingInfixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		left     int
		operator string
		right    int
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"5 <= 5", 5, "<=", 5},
		{"5 >= 5", 5, ">=", 5},
		{"5 >= 5", 5, ">=", 5},
	}
	for _, tt := range prefixTests {
		p := NewParser(NewLexer(tt.input))
		prog := p.ParseProgram()
		checkParserErrors(t, p)
		if len(prog.Statements) != 1 {
			t.Fatalf("expected 1 program statments. got=%d", len(prog.Statements))
		}
		stmt, ok := prog.Statements[0].(*ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not expression statement. got=%T", prog.Statements[0])
		}
		testInfixExpression(t, stmt.Expression, tt.left, tt.operator, tt.right)
	}
}

func testInfixExpression(t *testing.T, exp Expression, left any, op string, right any) bool {
	e, ok := exp.(*InfixExpression)
	if !ok {
		t.Fatalf("stmt is not infix expression. got=%T", exp)
	}
	if e.Operator != op {
		t.Fatalf("exp.operator is not '%s'. got='%s'",
			op, e.Operator)
	}
	if !testLiteralExpression(t, e.Left, left) {
		return false
	}
	if !testLiteralExpression(t, e.Right, right) {
		return false
	}
	return true
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"- a * b", "((-a)*b)"},
		{"! - a", "(!(-a))"},
		{"a + b + c", "((a+b)+c)"},
		{"a + b - c", "((a+b)-c)"},
		{"a * b * c", "((a*b)*c)"},
		{"a + b / c", "(a+(b/c))"},
		{"a + b * c + d / e - f", "(((a+(b*c))+(d/e))-f)"},
		{"3 + 4; -5 * 5", "(3+4)((-5)*5)"},
		{"5 > 4 != 3 < 4", "((5>4)!=(3<4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3+(4*5))==((3*1)+(4*5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3>5)==false)"},
		{"3 < 5 == true", "((3<5)==true)"},
		{"1 + (2 + 3) + 4", "((1+(2+3))+4)"},
		{"(5 + 5) * 2", "((5+5)*2)"},
		{"2 / (5 + 5)", "(2/(5+5))"},
		{"-(5 + 5)", "(-(5+5))"},
		{"!(true == true)", "(!(true==true))"},
		{"a + add(b * c) + d", "((a+add((b*c)))+d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a,b,1,(2*3),(4+5),add(6,(7*8)))"},
		{"add(a + b + c * d / f + g)",
			"add((((a+b)+((c*d)/f))+g))"},
	}
	for _, tt := range tests {
		p := NewParser(NewLexer(tt.input))
		prog := p.ParseProgram()
		checkParserErrors(t, p)
		actual := prog.String()
		if actual != tt.expected {
			t.Errorf("expected='%q', got='%q'", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input, expstr string
		expected      bool
	}{
		{"true;", "true", true},
		{"false;", "false", false},
	}
	for _, tt := range tests {
		p := NewParser(NewLexer(tt.input))
		prog := p.ParseProgram()
		checkParserErrors(t, p)
		if len(prog.Statements) != 1 {
			t.Fatalf("expected 1 program statments. got=%d", len(prog.Statements))
		}
		stmt, ok := prog.Statements[0].(*ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not expression statement. got=%T", prog.Statements[0])
		}
		testBooleanExpression(t, stmt.Expression, tt.expected)
	}
}

func testBooleanExpression(t *testing.T, stmt Expression, v bool) bool {
	b, ok := stmt.(*BooleanLiteral)
	if !ok {
		t.Fatalf("expression is not num literal. got=%T", stmt)
		return false
	}
	if b.Value != v {
		t.Errorf("literal value is not %t. got=%t", v, b.Value)
		return false
	}
	expstr := fmt.Sprint(v)
	if b.Literal != expstr {
		t.Errorf("ident.Literal not %s. got=%s", expstr, b.Literal)
		return false
	}
	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	p := NewParser(NewLexer(input))
	prog := p.ParseProgram()
	checkParserErrors(t, p)
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement. got=%d", len(prog.Statements))
	}
	stmt, ok := prog.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not expression. got=%T", prog.Statements[0])
	}
	exp, ok := stmt.Expression.(*IfExpression)
	if !ok {
		t.Fatalf("%s is not if expression. got=%T", stmt.Expression.String(), stmt.Expression)
	}
	testInfixExpression(t, exp.Condition, "x", "<", "y")
	if len(exp.Then.Statements) != 1 {
		t.Fatalf("then-leaf is not 1 statement. got=%d", len(exp.Then.Statements))
	}
	then, ok := exp.Then.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("'%s' is not expression. got=%T",
			exp.Then.Statements[0], exp.Then.Statements[0])
	}
	testIdentifier(t, then.Expression, "x")
	if exp.Else != nil {
		t.Errorf("else-leaf is not nil. got=%+v", exp.Else)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	p := NewParser(NewLexer(input))
	prog := p.ParseProgram()
	checkParserErrors(t, p)
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement. got=%d", len(prog.Statements))
	}
	stmt, ok := prog.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not expression. got=%T", prog.Statements[0])
	}
	exp, ok := stmt.Expression.(*IfExpression)
	if !ok {
		t.Fatalf("%s is not if expression. got=%T", stmt.Expression.String(), stmt.Expression)
	}
	testInfixExpression(t, exp.Condition, "x", "<", "y")
	if len(exp.Then.Statements) != 1 {
		t.Fatalf("then-leaf is not 1 statement. got=%d", len(exp.Then.Statements))
	}
	then, ok := exp.Then.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("'%s' is not expression. got=%T",
			exp.Then.Statements[0], exp.Then.Statements[0])
	}
	testIdentifier(t, then.Expression, "x")
	if exp.Else == nil {
		t.Error("else-leaf is nil.")
	}
	elseLeaf, ok := exp.Else.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("'%s' is not expression. got=%T",
			exp.Else.Statements[0], exp.Else.Statements[0])
	}
	testIdentifier(t, elseLeaf.Expression, "y")
}

func TestFuncLiteral(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	p := NewParser(NewLexer(input))
	prog := p.ParseProgram()
	checkParserErrors(t, p)
	if len(prog.Statements) != 1 {
		t.Fatalf("program doesn't have 1 statement. got=%d", len(prog.Statements))
	}
	stmt, ok := prog.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("'%s' is not expression stmt. got=%T", prog.Statements[0],
			prog.Statements[0])
	}
	fn, ok := stmt.Expression.(*FuncLiteral)
	if !ok {
		t.Fatalf("'%s' is not func literal. got=%T", stmt.Expression.String(),
			stmt.Expression)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("parametes expected 2. got=%d, %q", len(fn.Parameters), fn.Parameters)
	}
	testLiteralExpression(t, fn.Parameters[0], "x")
	testLiteralExpression(t, fn.Parameters[1], "y")

	if len(fn.Body.Statements) != 1 {
		t.Fatalf("fn body stmt expected 1. got=%T", len(fn.Body.Statements))
	}
	body, ok := fn.Body.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("'%s' is not expr stmt. got=%T",
			fn.Body.Statements[0],
			fn.Body.Statements[0],
		)
	}
	testInfixExpression(t, body.Expression, "x", "+", "y")

}

func TestFuncParamParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"fn(){};", []string{}},
		{"fn(x){};", []string{"x"}},
		{"fn(x,   y, z){};", []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		p := NewParser(NewLexer(tt.input))
		prog := p.ParseProgram()
		checkParserErrors(t, p)
		fn := prog.Statements[0].(*ExpressionStatement).
			Expression.(*FuncLiteral)
		if len(fn.Parameters) != len(tt.expected) {
			t.Errorf("length param expected %d. got=%d",
				len(tt.expected), len(fn.Parameters))
		}
		for i, id := range tt.expected {
			testLiteralExpression(t, fn.Parameters[i], id)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 *3, 4+5)"
	p := NewParser(NewLexer(input))
	prog := p.ParseProgram()
	checkParserErrors(t, p)
	if len(prog.Statements) != 1 {
		t.Fatalf("prog stmt expected 1. got=%d", len(prog.Statements))
	}
	stmt, ok := prog.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("'%s' is not call expression stmt. got=%T",
			prog.Statements[0], prog.Statements[0])
	}
	exp, ok := stmt.Expression.(*CallExpression)
	if !ok {
		t.Fatalf("'%s' is not call expr. got=%T", stmt.Expression, stmt.Expression)
	}
	testIdentifier(t, exp.Func, "add")
	if len(exp.Args) != 3 {
		t.Fatalf("expr arg expect 3. got=%T", len(exp.Args))
	}
	testLiteralExpression(t, exp.Args[0], 1)
	testInfixExpression(t, exp.Args[1], 2, "*", 3)
	testInfixExpression(t, exp.Args[2], 4, "+", 5)
}

func TestStringLiteral(t *testing.T) {
	input := `"hello world";`
	p := NewParser(NewLexer(input))
	prog := p.ParseProgram()
	checkParserErrors(t, p)
	if len(prog.Statements) != 1 {
		t.Fatalf("prog stmt expected 1. got=%d", len(prog.Statements))
	}
	stmt, ok := prog.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("'%s' is not string literal. got=%T",
			prog.Statements[0], prog.Statements[0])
	}
	str, ok := stmt.Expression.(*StringLiteral)
	if !ok {
		t.Fatalf("'%s' is not string literal. got=%T", stmt.Expression, stmt.Expression)
	}
	if str.Value != "hello world" {
		t.Errorf("expected hello world. got=%q", str.Value)
	}
}

func TestSliceLiteral(t *testing.T) {
	input := `[1, 2, hehehello, "aaa", true];`
	p := NewParser(NewLexer(input))
	prog := p.ParseProgram()
	checkParserErrors(t, p)
	if len(prog.Statements) != 1 {
		t.Fatalf("prog stmt expected 1. got=%d", len(prog.Statements))
	}
	stmt, ok := prog.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("'%s' is not expression. got=%T",
			prog.Statements[0], prog.Statements[0])
	}
	slc, ok := stmt.Expression.(*Slices)
	if !ok {
		t.Fatalf("'%s' is not slice. got=%T", stmt.Expression, stmt.Expression)
	}
	if len(slc.Elements) != 5 {
		t.Errorf("expected len 5. got=%q", len(slc.Elements))
	}
}
