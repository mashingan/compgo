package interp

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAst(t *testing.T) {
	const (
		mv = "myvar最強"
		nv = "zaワールド😊"
	)
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: Token{Let, "let"},
				Name: &Identifier{
					Token: Token{Ident, mv},
					Value: mv,
				},
				Value: &Identifier{
					Token: Token{Ident, nv},
					Value: nv,
				},
			},
		},
	}
	if program.String() != fmt.Sprintf("let %s = %s;", mv, nv) {
		t.Errorf("program string incorrect, got=%q", program.String())
	}

}

func TestModify(t *testing.T) {
	one := func() Expression { return &IntLiteral{Value: 1} }
	two := func() Expression { return &IntLiteral{Value: 2} }
	one2Two := func(node Node) Node {
		i, ok := node.(*IntLiteral)
		if !ok {
			return node
		}
		if i.Value != 1 {
			return node
		}
		i.Value = 2
		return i
	}
	tests := []struct{ input, expect Node }{
		{one(), two()},
		{&Program{
			Statements: []Statement{
				&ExpressionStatement{Expression: one()},
			},
		}, &Program{
			Statements: []Statement{
				&ExpressionStatement{Expression: two()},
			},
		},
		},
		{
			&InfixExpression{Left: one(), Operator: "+", Right: two()},
			&InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			&InfixExpression{Left: two(), Operator: "+", Right: one()},
			&InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			&PrefixExpression{Operator: "-", Right: one()},
			&PrefixExpression{Operator: "-", Right: two()},
		},
		{
			&CallIndex{Left: one(), Index: one()},
			&CallIndex{Left: two(), Index: two()},
		},
		{
			&IfExpression{
				Condition: one(),
				Then: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
				Else: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
			},
			&IfExpression{
				Condition: two(),
				Then: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
				Else: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
			},
		},
		{
			&ReturnStatement{Value: one()},
			&ReturnStatement{Value: two()},
		},
		{
			&LetStatement{Value: one()},
			&LetStatement{Value: two()},
		},
		{
			&FuncLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
			},
			&FuncLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
			},
		},
		{
			&Slices{Elements: []Expression{one(), one()}},
			&Slices{Elements: []Expression{two(), two()}},
		},
	}
	for _, tt := range tests {
		modified := Modify(tt.input, one2Two)
		eq := reflect.DeepEqual(modified, tt.expect)
		if !eq {
			t.Errorf("not equal. got=%#v, want=%#v",
				modified, tt.expect)
		}
	}
	hashLit := &HashLiteral{
		Pairs: map[Expression]Expression{
			one(): one(),
			one(): one(),
		},
	}
	Modify(hashLit, one2Two)
	for k, v := range hashLit.Pairs {
		kk, _ := k.(*IntLiteral)
		if kk.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, kk.Value)
		}
		vv, _ := v.(*IntLiteral)
		if vv.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, vv.Value)
		}
	}

}
