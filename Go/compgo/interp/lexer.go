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
	"=":      Assign,
	"+":      Plus,
	"(":      Lparen,
	")":      Rparen,
	"{":      Lbrace,
	"}":      Rbrace,
	",":      Comma,
	";":      Semicolon,
	"let":    Let,
	"fn":     Fn,
	"!":      Bang,
	"*":      Star,
	"/":      Slash,
	">":      Gt,
	"<":      Lt,
	"-":      Minus,
	"return": Return,
	"if":     If,
	"else":   Else,
	"true":   True,
	"false":  False,
	"!=":     Neq,
	">=":     Gte,
	"<=":     Lte,
	"==":     Eq,
	"\"":     Str,
	"[":      Lbracket,
	"]":      Rbracket,
}

func (l *Lexer) skipWhitespaces() {
	for r, size := utf8.DecodeRune(l.inputUtf8); unicode.IsSpace(r) && len(l.inputUtf8) > 0; r, size = utf8.DecodeRune(l.inputUtf8) {
		l.position++
		l.inputUtf8 = l.inputUtf8[size:]
	}
}

func (l *Lexer) getUntilSpaceOrOperator(p *[]byte) {
	for r, size := utf8.DecodeRune(l.inputUtf8); !unicode.IsSpace(r) && len(l.inputUtf8) > 0; r, size = utf8.DecodeRune(l.inputUtf8) {
		if _, ok := mapTokenLexer[string(r)]; ok {
			// in case of combined operator
			pp := utf8.AppendRune(*p, r)
			if _, ok := mapTokenLexer[string(pp)]; ok {
				*p = pp
				return
			}
			return
		}
		l.position++
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

func (l *Lexer) getCombined(t TokenType, r rune) Token {
	switch string(r) {
	case "=", "<", ">", "!":
		rr := utf8.AppendRune(nil, r)
		r, size := utf8.DecodeRune(l.inputUtf8)
		rr = utf8.AppendRune(rr, r)
		tt, ok := mapTokenLexer[string(rr)]
		if ok {
			l.position++
			l.inputUtf8 = l.inputUtf8[size:]
			return Token{tt, string(rr)}
		}
	case "\"":
		return l.readString()
	}
	return Token{t, string(r)}
}

func (l *Lexer) getToken() Token {
	l.skipWhitespaces()
	if len(l.inputUtf8) <= 0 {
		return Token{Eof, ""}
	}
	r, size := utf8.DecodeRune(l.inputUtf8)
	l.inputUtf8 = l.inputUtf8[size:]
	t, ok := mapTokenLexer[string(r)]
	l.position++
	if !ok {
		// means it's not a single rune
		return l.tokenize(r)
	}
	return l.getCombined(t, r)
}

func (l *Lexer) NextToken() Token {
	if len(l.inputUtf8) <= 0 {
		return Token{Eof, ""}
	}
	return l.getToken()
}

func (l *Lexer) readString() Token {
	rr := []byte{}
	escaped := false
	for {
		r, size := utf8.DecodeRune(l.inputUtf8)
		l.inputUtf8 = l.inputUtf8[size:]
		rs := string(r)
		if rs == "\"" && !escaped {
			break
		}
		if rs == "\\" && !escaped {
			escaped = true
			continue
		}
		if escaped {
			rr = appendEscape(rr, rs)
			escaped = false
			continue
		}
		rr = utf8.AppendRune(rr, r)
		l.position++
	}
	r, sz := utf8.DecodeRune(l.inputUtf8)
	if string(r) == "\"" {
		l.inputUtf8 = l.inputUtf8[sz:]
	}
	return Token{Str, string(rr)}
}

func appendEscape(rr []byte, rs string) []byte {
	switch rs {
	case "n":
		rr = append(rr, '\n')
	case "r":
		rr = append(rr, '\r')
	case "a":
		rr = append(rr, '\a')
	case "t":
		rr = append(rr, '\t')
	case "b":
		rr = append(rr, '\b')
	case "\\":
		rr = append(rr, '\\')
	case "\"":
		rr = append(rr, '"')
	}
	return rr
}
