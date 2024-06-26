package ast_test

import (
	"testing"

	"github.com/leapkit/leapkit/core/internal/plush/ast"
	"github.com/leapkit/leapkit/core/internal/plush/token"
	"github.com/stretchr/testify/require"
)

func Test_Program_String(t *testing.T) {
	r := require.New(t)
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				TokenAble: ast.TokenAble{token.Token{Type: token.LET, Literal: "let"}},
				Name: &ast.Identifier{
					TokenAble: ast.TokenAble{token.Token{Type: token.IDENT, Literal: "myVar"}},
					Value:     "myVar",
				},
				Value: &ast.Identifier{
					TokenAble: ast.TokenAble{token.Token{Type: token.IDENT, Literal: "anotherVar"}},
					Value:     "anotherVar",
				},
			},
		},
	}

	r.Equal("let myVar = anotherVar;", program.String())
}
