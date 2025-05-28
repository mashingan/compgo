package interp

import (
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
