package interp

import (
	"unicode"
	"unicode/utf8"
)

type pos struct {
	column  uint
	line    uint
	bytecol uint
}

type Lexer struct {
	inputUtf8    []byte
	inputStr     string
	position     int
	readPosition int
	ch           rune
	pos
}

func NewLexer(input string) *Lexer {
	return &Lexer{[]byte(input), input, 0, 0, 0, pos{line: 1}}
}

func (l *Lexer) newline() {
	l.pos.column = 0
	l.pos.bytecol++
	l.pos.line++
}

func (l *Lexer) forward(sz uint) {
	l.pos.column++
	l.pos.bytecol += sz
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
	":":      Colon,
	"macro":  Macro,
}

func (l *Lexer) skipWhitespaces() {
	for r, size := utf8.DecodeRune(l.inputUtf8); unicode.IsSpace(r) && len(l.inputUtf8) > 0; r, size = utf8.DecodeRune(l.inputUtf8) {
		if string(r) == "\n" {
			l.newline()
		} else {
			l.forward(uint(size))
		}
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
		l.forward(uint(size))
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
	lpos := l.pos
	lpos.column = max(lpos.column-uint(utf8.RuneCount(buf))+1, 0)
	if isNumber {
		return Token{Int, bstr, lpos}
	}
	t, ok := mapTokenLexer[bstr]
	if ok {
		return Token{t, bstr, lpos}
	}
	return Token{Ident, bstr, lpos}
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
			l.forward(uint(size))
			l.inputUtf8 = l.inputUtf8[size:]
			lpos := l.pos
			lpos.column -= 2 + 1
			return Token{tt, string(rr), l.pos}
		}
	case "\"":
		return l.readString()
	}
	lpos := l.pos
	lpos.column -= 1
	return Token{t, string(r), l.pos}
}

func (l *Lexer) getToken() Token {
	l.skipWhitespaces()
	if len(l.inputUtf8) <= 0 {
		return Token{Eof, "", l.pos}
	}
	r, size := utf8.DecodeRune(l.inputUtf8)
	l.forward(uint(size))
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
		return Token{Eof, "", l.pos}
	}
	return l.getToken()
}

func (l *Lexer) readString() Token {
	rr := []byte{}
	escaped := false
	forwardTimes := 0
	for {
		r, size := utf8.DecodeRune(l.inputUtf8)
		l.forward(uint(size))
		forwardTimes++
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
		l.forward(uint(sz))
		forwardTimes++
	}
	lpos := l.pos
	lpos.column -= uint(forwardTimes) + 1
	return Token{Str, string(rr), lpos}
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
