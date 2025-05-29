package interp

import (
	"fmt"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
	let x一 = 5;
	let y二 = 10;
	let foobar = 838383;
	`
	// 	input := `
	// let x一 5;
	// let = 10;
	// let 838383;
	// `
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram return nil")
	}
	checkParserErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d",
			len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x一"}, {"y二"}, {"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		testLetStatement(t, stmt, tt.expectedIdentifier)
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
	ident, ok := stmt.Expression.(*Identifier)
	if !ok {
		t.Fatalf("expression is not identfier. got=%T", stmt.Expression)
	}
	if ident.Value != input {
		t.Errorf("ident.Value not %s. got=%s", input, ident.Value)
	}
	if ident.Literal != input {
		t.Errorf("ident.Literal not %s. got=%s", input, ident.Literal)
	}
}

func TestIntLiteralExpression(t *testing.T) {
	input := "5;"
	expstr := "5"
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
	numint, ok := stmt.Expression.(*IntLiteral)
	if !ok {
		t.Fatalf("expression is not num literal. got=%T", stmt.Expression)
	}
	if numint.Value != expected {
		t.Errorf("literal value is not %d. got=%d", expected, numint.Value)
	}
	if numint.Literal != expstr {
		t.Errorf("ident.Literal not %s. got=%s", expstr, numint.Literal)
	}
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
		exp, ok := stmt.Expression.(*InfixExpression)
		if !ok {
			t.Fatalf("stmt is not infix expression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.operator is not '%s'. got='%s'",
				tt.operator, exp.Operator)
		}
		testIntLiteral(t, exp.Left, tt.left)
		if exp.Operator != tt.operator {
			t.Fatalf("exp operator is not '%s'. got='%s'",
				tt.operator, exp.Operator)
		}
		testIntLiteral(t, exp.Right, tt.right)
	}
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
