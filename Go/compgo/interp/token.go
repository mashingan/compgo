package interp

type TokenType uint

type Token struct {
	Type    TokenType
	Literal string
}

const (
	Illegal TokenType = iota
	Eof
	Ident
	Int
	Assign
	Plus
	Minus
	Star
	Slash
	Bang
	Comma
	Semicolon
	Lparen
	Rparen
	Lbrace
	Rbrace
	Function
	Let
	Lt
	Gt
	Return
	If
	Else
	True
	False
	Eq
	Gte
	Lte
	Neq
)
