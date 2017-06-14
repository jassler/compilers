package lexer

import "testing"
import "github.com/simplang/token"

func TestNextToken(t *testing.T) {
	input := `let a = 1 and
	loopy = a+-1
in
	loopy
end`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "1"},
		{token.AND, "and"},
		{token.IDENT, "loopy"},
		{token.ASSIGN, "="},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.INT, "1"},
		{token.IN, "in"},
		{token.IDENT, "loopy"},
		{token.END, "end"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
