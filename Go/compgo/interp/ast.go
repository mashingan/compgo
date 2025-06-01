package interp

import (
	"fmt"
	"strings"
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
	var out strings.Builder
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

type IntLiteral struct {
	Token
	Value int
}

func (i *IntLiteral) expressionNode()      {}
func (i *IntLiteral) TokenLiteral() string { return i.Literal }
func (i *IntLiteral) String() string       { return i.Literal }

type PrefixExpression struct {
	Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Literal }
func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Operator, p.Right.String())
}

type InfixExpression struct {
	Token
	Operator    string
	Left, Right Expression
}

func (p *InfixExpression) expressionNode()      {}
func (p *InfixExpression) TokenLiteral() string { return p.Literal }
func (p *InfixExpression) String() string {
	return fmt.Sprintf("(%s%s%s)", p.Left.String(), p.Operator, p.Right.String())
}

type BooleanLiteral struct {
	Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Literal }
func (b *BooleanLiteral) String() string       { return b.Literal }

type IfExpression struct {
	Token
	Condition  Expression
	Then, Else *BlockStatement
}

func (i *IfExpression) expressionNode()      {}
func (i *IfExpression) TokenLiteral() string { return i.Literal }
func (i *IfExpression) String() string {
	elseLeaf := ""
	if i.Else != nil {
		elseLeaf = fmt.Sprintf("else %s", i.Else.String())
	}
	return fmt.Sprintf("if %s %s%s", i.Condition.String(), i.Then.String(), elseLeaf)
}

type BlockStatement struct {
	Token
	Statements []Statement
}

func (b *BlockStatement) expressionNode()      {}
func (b *BlockStatement) TokenLiteral() string { return b.Literal }
func (b *BlockStatement) String() string {
	var sb strings.Builder
	for _, s := range b.Statements {
		sb.WriteString(s.String())
	}
	return sb.String()
}

type FuncLiteral struct {
	Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FuncLiteral) expressionNode()      {}
func (f *FuncLiteral) TokenLiteral() string { return f.Literal }
func (f *FuncLiteral) String() string {
	var sb strings.Builder
	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.String()
	}
	sb.WriteString(f.Literal)
	sb.WriteString(fmt.Sprintf("(%s)", strings.Join(params, ",")))
	sb.WriteString(f.Body.String())
	return sb.String()
}

type CallExpression struct {
	Token
	Func Expression
	Args []Expression
}

func (c *CallExpression) expressionNode()      {}
func (c *CallExpression) TokenLiteral() string { return c.Literal }
func (c *CallExpression) String() string {
	params := make([]string, len(c.Args))
	for i, p := range c.Args {
		params[i] = p.String()
	}
	return fmt.Sprintf("%s(%s)", c.Func.String(), strings.Join(params, ","))
}

type StringLiteral struct {
	Token
	Value string
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Literal }
func (s *StringLiteral) String() string       { return s.Literal }

type Slices struct {
	Token
	Elements []Expression
}

func (s *Slices) expressionNode()      {}
func (s *Slices) TokenLiteral() string { return s.Literal }
func (s *Slices) String() string {
	res := make([]string, len(s.Elements))
	for i, e := range s.Elements {
		res[i] = e.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(res, ","))
}

type CallIndex struct {
	Token
	Left, Index Expression
}

func (c *CallIndex) expressionNode()      {}
func (c *CallIndex) TokenLiteral() string { return c.Literal }
func (c *CallIndex) String() string {
	return fmt.Sprintf("%s[%s]", c.Left, c.Index)
}
