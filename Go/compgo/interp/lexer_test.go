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

	input = `let five異能 = 5;
let ten世界 = 10;

let add特異点 = fn(x, y) {
	x + y;
};

let result = add特異点(five, ten);
!-/*5;
5 < 10 > 5;
`

	tests = []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{Let, "let"},
		{Ident, "five異能"},
		{Assign, "="},
		{Int, "5"},
		{Semicolon, ";"},
		{Let, "let"},
		{Ident, "ten世界"},
		{Assign, "="},
		{Int, "10"},
		{Semicolon, ";"},
		{Let, "let"},
		{Ident, "add特異点"},
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
		{Ident, "add特異点"},
		{Lparen, "("},
		{Ident, "five"},
		{Comma, ","},
		{Ident, "ten"},
		{Rparen, ")"},
		{Semicolon, ";"},
		{Bang, "!"},
		{Minus, "-"},
		{Slash, "/"},
		{Star, "*"},
		{Int, "5"},
		{Semicolon, ";"},
		{Int, "5"},
		{Lt, "<"},
		{Int, "10"},
		{Gt, ">"},
		{Int, "5"},
		{Semicolon, ";"},
		{Eof, ""},
	}
	l = NewLexer(input)
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
}
