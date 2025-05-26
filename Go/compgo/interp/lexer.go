package interp

import "unicode/utf8"

type Lexer struct {
	inputUtf8    []byte
	inputStr     string
	position     int
	readPosition int
	ch           rune
}

func NewLexer(input string) *Lexer {
	return &Lexer{[]byte(input), input, 0, 0, 0}
}

var mapTokenLexer = map[string]TokenType{
	"=": Assign,
	"+": Plus,
	"(": Lparen,
	")": Rparen,
	"{": Lbrace,
	"}": Rbrace,
	",": Comma,
	";": Semicolon,
}

func (l *Lexer) NextToken() Token {
	if len(l.inputUtf8) <= 0 {
		return Token{Eof, ""}
	}
	r, size := utf8.DecodeRune(l.inputUtf8)
	l.readPosition += size
	l.inputUtf8 = l.inputUtf8[size:]
	tch := string(r)
	l.position++
	t, ok := mapTokenLexer[tch]
	if !ok {
		return Token{Illegal, ""}
	}
	return Token{t, tch}
}
