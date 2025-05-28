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

func (t TokenType) String() string {
	s := mapTokenDisplay[t]
	return s
}

var mapTokenDisplay = map[TokenType]string{
	Assign:    "=",
	Plus:      "+",
	Lparen:    "(",
	Rparen:    ")",
	Lbrace:    "{",
	Rbrace:    "}",
	Comma:     ",",
	Semicolon: ";",
	Let:       "let",
	Function:  "fn",
	Bang:      "!",
	Star:      "*",
	Slash:     "/",
	Gt:        ">",
	Lt:        "<",
	Minus:     "-",
	Return:    "return",
	If:        "if",
	Else:      "else",
	True:      "true",
	False:     "false",
	Neq:       "!=",
	Gte:       ">=",
	Lte:       "<=",
	Eq:        "==",
	Ident:     "Ident",
	Int:       "Int",
}
