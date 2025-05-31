package interp

import "testing"

func TestNextToken_1(t *testing.T) {
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
}

func TestNextToken_2(t *testing.T) {
	input := `let five異能 = 5;
let ten世界 = 10;

let add特異点 = fn(x, y) {
	x + y;
};

let result = add特異点(five, ten);
`
	tests := []struct {
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
		{Fn, "fn"},
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
}
func TestNextToken_3(t *testing.T) {
	input := `!-/*5;
	5 < 10 > 5;
	`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{

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
}

func TestNextToken_4(t *testing.T) {

	input := `if (5 < 10) {
		return true;
	} else if (10 > 5) {
		return false;
	}`
	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{If, "if"},
		{Lparen, "("},
		{Int, "5"},
		{Lt, "<"},
		{Int, "10"},
		{Rparen, ")"},
		{Lbrace, "{"},
		{Return, "return"},
		{True, "true"},
		{Semicolon, ";"},
		{Rbrace, "}"},
		{Else, "else"},
		{If, "if"},
		{Lparen, "("},
		{Int, "10"},
		{Gt, ">"},
		{Int, "5"},
		{Rparen, ")"},
		{Lbrace, "{"},
		{Return, "return"},
		{False, "false"},
		{Semicolon, ";"},
		{Rbrace, "}"},
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

}

func TestNextToken_5(t *testing.T) {

	input := `5 == 10; 10 != 9;
5 <= 9;
10 >= 10;
`
	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{Int, "5"}, {Eq, "=="}, {Int, "10"}, {Semicolon, ";"},
		{Int, "10"}, {Neq, "!="}, {Int, "9"}, {Semicolon, ";"},
		{Int, "5"}, {Lte, "<="}, {Int, "9"}, {Semicolon, ";"},
		{Int, "10"}, {Gte, ">="}, {Int, "10"}, {Semicolon, ";"},
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

}
