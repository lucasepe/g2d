package lexer

import (
	"testing"

	"github.com/lucasepe/g2d/token"
)

func TestNextToken(t *testing.T) {
	input := `#!./g2d
five := 5;
ten := 10;
fl := 8.88

add := fn(x, y) {
  x + y;
};

# this is a comment
result := add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
    return true;
} else {
    return false;
}

// this is another comment
10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2]
{"foo": "bar"}
d.foo
!&&||
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.COMMENT, "!./g2d"},
		{token.IDENT, "five"},
		{token.BIND, ":="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "ten"},
		{token.BIND, ":="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "fl"},
		{token.BIND, ":="},
		{token.FLOAT, "8.88"},
		{token.IDENT, "add"},
		{token.BIND, ":="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.COMMENT, " this is a comment"},
		{token.IDENT, "result"},
		{token.BIND, ":="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.NOT, "!"},
		{token.MINUS, "-"},
		{token.DIVIDE, "/"},
		{token.MULTIPLY, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.COMMENT, " this is another comment"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NEQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.IDENT, "d"},
		{token.DOT, "."},
		{token.IDENT, "foo"},
		{token.NOT, "!"},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, test := range tests {
		token := lexer.NextToken()

		if token.Type != test.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q",
				i, test.expectedType, token.Type)
		}

		if token.Literal != test.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, test.expectedLiteral, token.Literal)
		}
	}
}

func TestStringEscapes(t *testing.T) {
	input := `#!./g2d
a := "\"foo\""
b := "\x00\x0a\x7f"
c := "\r\n\t"
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.COMMENT, "!./g2d"},
		{token.IDENT, "a"},
		{token.BIND, ":="},
		{token.STRING, "\"foo\""},
		{token.IDENT, "b"},
		{token.BIND, ":="},
		{token.STRING, "\x00\n\u007f"},
		{token.IDENT, "c"},
		{token.BIND, ":="},
		{token.STRING, "\r\n\t"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, test := range tests {
		token := lexer.NextToken()

		if token.Type != test.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q",
				i, test.expectedType, token.Type)
		}

		if token.Literal != test.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, test.expectedLiteral, token.Literal)
		}
	}

}
