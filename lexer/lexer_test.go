package lexer

import (
	"testing"

	"github.com/tomocy/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := "=+(){},;"
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{
			expectedType:    token.ASSIGN,
			expectedLiteral: "=",
		},
		{
			expectedType:    token.PLUS,
			expectedLiteral: "+",
		},
		{
			expectedType:    token.LPAREN,
			expectedLiteral: "(",
		},
		{
			expectedType:    token.RPAREN,
			expectedLiteral: ")",
		},
		{
			expectedType:    token.LBRACE,
			expectedLiteral: "{",
		},
		{
			expectedType:    token.RBRACE,
			expectedLiteral: "}",
		},
		{
			expectedType:    token.COMMA,
			expectedLiteral: ",",
		},
		{
			expectedType:    token.SEMICOLON,
			expectedLiteral: ";",
		},
	}

	lexer := New(input)

	for i, test := range tests {
		nextToken := lexer.NextToken()
		if nextToken.Type != test.expectedType {
			t.Errorf("test[%d]: wrong type. expected %s, but got %s\n", i, test.expectedType, nextToken.Type)
		}
		if nextToken.Literal != test.expectedLiteral {
			t.Errorf("test[%d]: wrong literal. expected %s, but got %s\n", i, test.expectedLiteral, nextToken.Literal)
		}
	}
}
