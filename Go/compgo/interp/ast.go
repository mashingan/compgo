package interp

import (
	"bytes"
	"fmt"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type Identifier struct {
	Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Literal }
func (i *Identifier) String() string       { return i.Value }

type LetStatement struct {
	Token
	Name  *Identifier
	Value Expression
}

func (l *LetStatement) statementNode()       {}
func (l *LetStatement) TokenLiteral() string { return l.Literal }
func (l *LetStatement) String() string {
	val := ""
	if l.Value != nil {
		val = fmt.Sprintf(" = %s", l.Value.String())
	}
	return fmt.Sprintf("let %s%s;", l.Name, val)
}

type ReturnStatement struct {
	Token
	Value Expression
}

func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) TokenLiteral() string { return r.Literal }
func (r *ReturnStatement) String() string {
	val := ""
	if r.Value != nil {
		val = fmt.Sprintf(" %s", r.Value.String())
	}
	return fmt.Sprintf("return%s;", val)
}

type ExpressionStatement struct {
	Token
	Expression
}

func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Literal }
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}
