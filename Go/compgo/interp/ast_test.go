package interp

import (
	"fmt"
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
