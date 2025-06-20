package interp

import (
	"fmt"
	"strconv"
)

const (
	_ uint8 = iota
	Lowest
	Equals
	Lessgreater
	Sum
	Product
	Prefix
	Call
	Index
)

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)
type Parser struct {
	l                    *Lexer
	currToken, peekToken Token
	errors               []string

	prefixs map[TokenType]prefixParseFn
	infixs  map[TokenType]infixParseFn
}

var precedences = map[TokenType]uint8{
	Eq:       Equals,
	Neq:      Equals,
	Lt:       Lessgreater,
	Gt:       Lessgreater,
	Lte:      Lessgreater,
	Gte:      Lessgreater,
	Plus:     Sum,
	Minus:    Sum,
	Slash:    Product,
	Star:     Product,
	Lparen:   Call,
	Lbracket: Index,
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.prefixs = map[TokenType]prefixParseFn{}
	p.prefixs[Ident] = p.parseIdentifier
	p.prefixs[Int] = p.parseIntLiteral
	p.prefixs[Bang] = p.parsePrefixExpression
	p.prefixs[Minus] = p.parsePrefixExpression
	p.prefixs[True] = p.parseBoolean
	p.prefixs[False] = p.parseBoolean
	p.prefixs[Lparen] = p.parseGroupExpression
	p.prefixs[If] = p.parseIfExpression
	p.prefixs[Fn] = p.parseFuncLiteral
	p.prefixs[Str] = p.parseStrLiteral
	p.prefixs[Lbracket] = p.parseSlice
	p.prefixs[Lbrace] = p.parseHashMap
	p.prefixs[Macro] = p.parseMacroLiteral
	p.infixs = map[TokenType]infixParseFn{}
	p.infixs[Plus] = p.parseInfixExpression
	p.infixs[Minus] = p.parseInfixExpression
	p.infixs[Star] = p.parseInfixExpression
	p.infixs[Slash] = p.parseInfixExpression
	p.infixs[Gt] = p.parseInfixExpression
	p.infixs[Lt] = p.parseInfixExpression
	p.infixs[Eq] = p.parseInfixExpression
	p.infixs[Neq] = p.parseInfixExpression
	p.infixs[Gte] = p.parseInfixExpression
	p.infixs[Lte] = p.parseInfixExpression
	p.infixs[Lparen] = p.parseCallExpression
	p.infixs[Lbracket] = p.parseIndexing
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string { return p.errors }
func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token %s, got %s",
		t, p.peekToken.Literal)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}
	for p.currToken.Type != Eof {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.currToken.Type {
	case Let:
		return p.parseLetStatement()
	case Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) expectNext(t TokenType) bool {
	if p.peekToken.Type != t {
		p.peekError(t)
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.currToken}
	if !p.expectNext(Ident) {
		return nil
	}
	stmt.Name = &Identifier{Token: p.currToken, Value: p.currToken.Literal}
	if !p.expectNext(Assign) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(Lowest)
	if p.peekToken.Type == Semicolon {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.currToken}
	p.nextToken()
	stmt.Value = p.parseExpression(Lowest)
	if p.peekToken.Type == Semicolon {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.currToken}
	stmt.Expression = p.parseExpression(Lowest)
	if p.peekToken.Type == Semicolon {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence uint8) Expression {
	prefix := p.prefixs[p.currToken.Type]
	if prefix == nil {
		msg := fmt.Sprintf("no prefix parse function for %s found", p.currToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	left := prefix()
	for p.peekToken.Type != Semicolon && precedence < p.peekPrecedence() {
		infix := p.infixs[p.peekToken.Type]
		if infix == nil {
			return left
		}
		p.nextToken()
		left = infix(left)
	}
	return left
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntLiteral() Expression {
	lit := &IntLiteral{Token: p.currToken}
	value, err := strconv.Atoi(p.currToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("cannot parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpression() Expression {
	exp := &PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}
	p.nextToken()
	exp.Right = p.parseExpression(Prefix)
	return exp
}

func (p *Parser) peekPrecedence() uint8 {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return Lowest
}

func (p *Parser) currPrecedence() uint8 {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return Lowest
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	e := &InfixExpression{
		Left:     left,
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}
	pred := p.currPrecedence()
	p.nextToken()
	e.Right = p.parseExpression(pred)
	return e
}

func (p *Parser) parseBoolean() Expression {
	return &BooleanLiteral{p.currToken, p.currToken.Type == True}
}

func (p *Parser) parseGroupExpression() Expression {
	p.nextToken()
	exp := p.parseExpression(Lowest)
	if !p.expectNext(Rparen) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() Expression {
	ifexp := &IfExpression{Token: p.currToken}
	if !p.expectNext(Lparen) {
		return nil
	}
	p.nextToken()
	ifexp.Condition = p.parseExpression(Lowest)
	if !p.expectNext(Rparen) {
		return nil
	}
	if !p.expectNext(Lbrace) {
		return nil
	}
	ifexp.Then = p.parseBlockStatement()
	if p.peekToken.Type == Else {
		p.nextToken()
		if !p.expectNext(Lbrace) {
			return nil
		}
		ifexp.Else = p.parseBlockStatement()
	}
	return ifexp
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	b := &BlockStatement{Token: p.currToken}
	b.Statements = []Statement{}
	p.nextToken()
	for p.currToken.Type != Rbrace && p.currToken.Type != Eof {
		b.Statements = append(b.Statements, p.parseStatement())
		p.nextToken()
	}
	return b
}

func (p *Parser) parseFuncLiteral() Expression {
	fn := &FuncLiteral{Token: p.currToken}
	if !p.expectNext(Lparen) {
		return nil
	}
	p.nextToken()
	for p.currToken.Type != Rparen {
		fn.Parameters = append(fn.Parameters,
			&Identifier{p.currToken, p.currToken.Literal})
		p.nextToken()
		if p.currToken.Type == Comma {
			p.nextToken()
		}
	}
	if !p.expectNext(Lbrace) {
		return nil
	}
	fn.Body = p.parseBlockStatement()
	return fn
}

func (p *Parser) parseListUntil(tokenType TokenType) []Expression {
	lst := []Expression{}
	p.nextToken()
	for p.currToken.Type != tokenType {
		lst = append(lst, p.parseExpression(Lowest))
		p.nextToken()
		if p.currToken.Type == Comma {
			p.nextToken()
		}
	}
	return lst
}

func (p *Parser) parseCallExpression(fn Expression) Expression {
	ce := &CallExpression{Token: p.currToken, Func: fn}
	ce.Args = p.parseListUntil(Rparen)
	return ce
}

func (p *Parser) parseStrLiteral() Expression {
	return &StringLiteral{p.currToken, p.currToken.Literal}
}

func (p *Parser) parseSlice() Expression {
	e := &Slices{Token: p.currToken}
	e.Elements = p.parseListUntil(Rbracket)
	return e
}

func (p *Parser) parseIndexing(left Expression) Expression {
	c := &CallIndex{Token: p.currToken, Left: left}
	p.nextToken()
	c.Index = p.parseExpression(Lowest)
	p.nextToken()
	return c
}

func (p *Parser) parseHashMap() Expression {
	h := &HashLiteral{p.currToken, make(map[Expression]Expression)}
	p.nextToken()
	for p.currToken.Type != Rbrace {
		left := p.parseExpression(Lowest)
		if !p.expectNext(Colon) {
			p.peekError(Colon)
			return nil
		}
		p.nextToken()
		right := p.parseExpression(Lowest)
		p.nextToken()
		h.Pairs[left] = right
		if p.currToken.Type != Comma && p.currToken.Type != Rbrace {
			p.peekError(Comma)
			return nil
		}
		if p.currToken.Type != Rbrace {
			p.nextToken()
		}
	}
	return h
}

func (p *Parser) parseMacroLiteral() Expression {
	fn := p.parseFuncLiteral()
	ffn, _ := fn.(*FuncLiteral)
	return &MacroLiteral{Token: ffn.Token, Parameters: ffn.Parameters,
		Body: ffn.Body}
}
