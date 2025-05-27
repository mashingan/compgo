package interp

import (
	"unicode"
	"unicode/utf8"
)

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
	"=":   Assign,
	"+":   Plus,
	"(":   Lparen,
	")":   Rparen,
	"{":   Lbrace,
	"}":   Rbrace,
	",":   Comma,
	";":   Semicolon,
	"let": Let,
	"fn":  Function,
	"!":   Bang,
	"*":   Star,
	"/":   Slash,
	">":   Gt,
	"<":   Lt,
	"-":   Minus,
}

func (l *Lexer) skipWhitespaces() {
	for r, size := utf8.DecodeRune(l.inputUtf8); unicode.IsSpace(r) && len(l.inputUtf8) > 0; r, size = utf8.DecodeRune(l.inputUtf8) {
		l.inputUtf8 = l.inputUtf8[size:]
	}
}

func (l *Lexer) getUntilSpaceOrOperator(p *[]byte) {
	for r, size := utf8.DecodeRune(l.inputUtf8); !unicode.IsSpace(r) && len(l.inputUtf8) > 0; r, size = utf8.DecodeRune(l.inputUtf8) {
		if _, ok := mapTokenLexer[string(r)]; ok {
			// skip in case find (){},;+=
			return
		}
		l.inputUtf8 = l.inputUtf8[size:]
		*p = utf8.AppendRune(*p, r)
	}
}

func (l *Lexer) tokenize(r rune) Token {
	buf := utf8.AppendRune(nil, r)
	l.getUntilSpaceOrOperator(&buf)
	cursize := 0
	isNumber := true
	for rr, rsize := utf8.DecodeRune(buf[cursize:]); len(buf[cursize:]) > 0; rr, rsize = utf8.DecodeRune(buf[cursize:]) {
		if !unicode.IsNumber(rr) {
			isNumber = false
			break
		}
		cursize += rsize
	}
	bstr := string(buf)
	if isNumber {
		return Token{Int, bstr}
	}
	t, ok := mapTokenLexer[bstr]
	if ok {
		return Token{t, bstr}
	}
	return Token{Ident, bstr}
}

func (l *Lexer) getToken() Token {
	l.skipWhitespaces()
	if len(l.inputUtf8) <= 0 {
		return Token{Eof, ""}
	}
	r, size := utf8.DecodeRune(l.inputUtf8)
	l.inputUtf8 = l.inputUtf8[size:]
	t, ok := mapTokenLexer[string(r)]
	if !ok {
		// means it's not a single rune
		return l.tokenize(r)
	}
	return Token{t, string(r)}
}

func (l *Lexer) NextToken() Token {
	if len(l.inputUtf8) <= 0 {
		return Token{Eof, ""}
	}
	return l.getToken()
}
