package interp

import (
	"fmt"
)

type Parser struct {
	l                    *Lexer
	currToken, peekToken Token
	errors               []string
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
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
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
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
	}
	return nil
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
	for p.currToken.Type != Semicolon {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.currToken}
	p.nextToken()
	for p.currToken.Type != Semicolon {
		p.nextToken()
	}
	return stmt
}
