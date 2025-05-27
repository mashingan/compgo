package interp

import "testing"

func TestNextToken(t *testing.T) {
	input := `=+(){},;`
	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{Assign, "="},
		{Plus, "+"},
		{Lparen, "("},
		{Rparen, ")"},
		{Lbrace, "{"},
		{Rbrace, "}"},
		{Comma, ","},
		{Semicolon, ";"},
		{Eof, ""},
	}
	l := NewLexer(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got %q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - token literal wrong. expected=%q, got %q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}

	input = `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);
`

	tests = []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{Let, "let"},
		{Ident, "five"},
		{Assign, "="},
		{Int, "5"},
		{Semicolon, ";"},
		{Let, "let"},
		{Ident, "ten"},
		{Assign, "="},
		{Int, "10"},
		{Semicolon, ";"},
		{Let, "let"},
		{Ident, "add"},
		{Assign, "="},
		{Function, "fn"},
		{Lparen, "("},
		{Ident, "x"},
		{Comma, ","},
		{Ident, "y"},
		{Rparen, ")"},
		{Lbrace, "{"},
		{Ident, "x"},
		{Plus, "+"},
		{Ident, "y"},
		{Semicolon, ";"},
		{Rbrace, "}"},
		{Semicolon, ";"},
		{Let, "let"},
		{Ident, "result"},
		{Assign, "="},
		{Ident, "add"},
		{Lparen, "("},
		{Ident, "five"},
		{Comma, ","},
		{Ident, "ten"},
		{Rparen, ")"},
		{Semicolon, ";"},
		{Eof, ""},
	}
}
