package interp

import (
	"fmt"
	"strings"
	"sync"
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

type HashLiteral struct {
	Token
	Pairs map[Expression]Expression
}

func (h *HashLiteral) expressionNode()      {}
func (h *HashLiteral) TokenLiteral() string { return h.Literal }
func (h *HashLiteral) String() string {
	bd := make([]string, len(h.Pairs))
	idx := 0
	for k, v := range h.Pairs {
		bd[idx] = fmt.Sprintf("%s:%s", k, v)
		idx++
	}
	return fmt.Sprintf("{%s}", strings.Join(bd, ","))
}

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i] = Modify(statement, modifier).(Statement)
		}
	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)
	case *InfixExpression:
		var w sync.WaitGroup
		w.Add(2)
		go func(w *sync.WaitGroup) {
			defer w.Done()
			node.Left, _ = Modify(node.Left, modifier).(Expression)
		}(&w)
		go func(w *sync.WaitGroup) {
			defer w.Done()
			node.Right, _ = Modify(node.Right, modifier).(Expression)
		}(&w)
		w.Wait()
	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *CallIndex:
		var w sync.WaitGroup
		w.Add(2)
		go func(w *sync.WaitGroup) {
			defer w.Done()
			node.Left, _ = Modify(node.Left, modifier).(Expression)
		}(&w)
		go func(w *sync.WaitGroup) {
			defer w.Done()
			node.Index, _ = Modify(node.Index, modifier).(Expression)
		}(&w)
		w.Wait()
	case *IfExpression:
		var w sync.WaitGroup
		w.Add(2)
		go func(w *sync.WaitGroup) {
			defer w.Done()
			node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		}(&w)
		go func(w *sync.WaitGroup) {
			defer w.Done()
			node.Then, _ = Modify(node.Then, modifier).(*BlockStatement)
		}(&w)
		if node.Else != nil {
			w.Add(1)
			go func(w *sync.WaitGroup) {
				defer w.Done()
				node.Else, _ = Modify(node.Else, modifier).(*BlockStatement)
			}(&w)
		}
		w.Wait()
	case *BlockStatement:
		var w sync.WaitGroup
		for i := range node.Statements {
			w.Add(1)
			go func(i int, w *sync.WaitGroup) {
				defer w.Done()
				node.Statements[i], _ = Modify(node.Statements[i], modifier).(Statement)
			}(i, &w)
		}
		w.Wait()
	case *ReturnStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)
	case *LetStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)
	case *FuncLiteral:
		var w sync.WaitGroup
		w.Add(1)
		go func(w *sync.WaitGroup) {
			defer w.Done()
			node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)
		}(&w)
		for i := range node.Parameters {
			w.Add(1)
			go func(i int, w *sync.WaitGroup) {
				defer w.Done()
				node.Parameters[i], _ = Modify(node.Parameters[i], modifier).(*Identifier)
			}(i, &w)
		}
		w.Wait()
	case *Slices:
		var w sync.WaitGroup
		for i := range node.Elements {
			w.Add(1)
			go func(i int, w *sync.WaitGroup) {
				defer w.Done()
				node.Elements[i], _ = Modify(node.Elements[i], modifier).(Expression)
			}(i, &w)
		}
		w.Wait()
	case *HashLiteral:
		for k, v := range node.Pairs {
			kk, _ := Modify(k, modifier).(Expression)
			vv, _ := Modify(v, modifier).(Expression)
			delete(node.Pairs, k)
			node.Pairs[kk] = vv
		}
	}
	return modifier(node)
}

type MacroLiteral struct {
	Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (m *MacroLiteral) expressionNode()      {}
func (m *MacroLiteral) TokenLiteral() string { return m.Literal }
func (m *MacroLiteral) String() string {
	params := make([]string, len(m.Parameters))
	for i, p := range m.Parameters {
		params[i] = p.String()
	}
	return fmt.Sprintf("macro(%s)%s", strings.Join(params, ","), m.Body.String())
}

func ExpandMacros(prg Node, env *Environment) Node {
	return Modify(prg, func(n Node) Node {
		ce, ok := n.(*CallExpression)
		if !ok {
			return n
		}
		id, ok := ce.Func.(*Identifier)
		if !ok {
			return n
		}
		o, ok := env.Get(id.Value)
		if !ok {
			return n
		}
		macro, ok := o.(*MacroObj)
		if !ok {
			return n
		}
		newenv := NewEnvironmentFrame(macro.Env)
		for i, p := range macro.Parameters {
			newenv.Set(p.Value, &Quote{ce.Args[i]})
		}
		evl := Eval(macro.Body, newenv)
		quote, ok := evl.(*Quote)
		if !ok {
			return n
		}
		return quote.Node
	})
}
