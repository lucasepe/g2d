package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lucasepe/g2d/token"
)

func TestString(t *testing.T) {
	assert := assert.New(t)

	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Token: token.Token{Type: token.IDENT, Literal: "myVar"},
				Expression: &BindExpression{
					Token: token.Token{Type: token.BIND, Literal: ":="},
					Left: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "myVar"},
						Value: "myVar",
					},
					Value: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
						Value: "anotherVar",
					},
				},
			},
		},
	}

	assert.Equal("myVar:=anotherVar", program.String())
}
