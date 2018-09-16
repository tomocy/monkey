package lexer

import (
	"testing"

	"github.com/tomocy/monkey/token"
)

func TestNextToken(t *testing.T) {
	type expect struct {
		tokenType token.TokenType
		literal   string
	}
	tests := []struct {
		in      string
		expects []expect
	}{
		{
			`
			let five = 5;
			let ten = 10;
			let add = fn(x, y) {
				x + y;
			};
			let result = add(five, ten);
			`,
			[]expect{
				{token.Let, "let"}, {token.Ident, "five"}, {token.Assign, "="}, {token.Int, "5"}, {token.Semicolon, ";"},
				{token.Let, "let"}, {token.Ident, "ten"}, {token.Assign, "="}, {token.Int, "10"}, {token.Semicolon, ";"},
				{token.Let, "let"}, {token.Ident, "add"}, {token.Assign, "="},
				{token.Function, "fn"}, {token.LParen, "("}, {token.Ident, "x"}, {token.Comma, ","}, {token.Ident, "y"}, {token.RParen, ")"},
				{token.LBrace, "{"}, {token.Ident, "x"}, {token.Plus, "+"}, {token.Ident, "y"}, {token.Semicolon, ";"}, {token.RBrace, "}"},
				{token.Semicolon, ";"},
				{token.Let, "let"}, {token.Ident, "result"}, {token.Assign, "="},
				{token.Ident, "add"}, {token.LParen, "("}, {token.Ident, "five"}, {token.Comma, ","}, {token.Ident, "ten"}, {token.RParen, ")"},
				{token.Semicolon, ";"},
				{token.EOF, ""},
			},
		},
		{
			`
			!-/*5;
			5 < 10 > 5;
			10 == 10;
			10 != 10;
			`,
			[]expect{
				{token.Bang, "!"}, {token.Minus, "-"}, {token.Slash, "/"}, {token.Asterrisk, "*"}, {token.Int, "5"}, {token.Semicolon, ";"},
				{token.Int, "5"}, {token.LessThan, "<"}, {token.Int, "10"}, {token.GreaterThan, ">"}, {token.Int, "5"}, {token.Semicolon, ";"},
				{token.Int, "10"}, {token.Equal, "=="}, {token.Int, "10"}, {token.Semicolon, ";"},
				{token.Int, "10"}, {token.NotEqual, "!="}, {token.Int, "10"}, {token.Semicolon, ";"},
				{token.EOF, ""},
			},
		},
		{
			`
			if (5 < 10) {
				return true;
			} else {
				return false;
			}
			`,
			[]expect{
				{token.If, "if"}, {token.LParen, "("}, {token.Int, "5"}, {token.LessThan, "<"}, {token.Int, "10"}, {token.RParen, ")"},
				{token.LBrace, "{"}, {token.Return, "return"}, {token.True, "true"}, {token.Semicolon, ";"}, {token.RBrace, "}"},
				{token.Else, "else"}, {token.LBrace, "{"}, {token.Return, "return"}, {token.False, "false"}, {token.Semicolon, ";"}, {token.RBrace, "}"},
				{token.EOF, ""},
			},
		},
		{
			`
			"foobar";
			"foo bar";
			`,
			[]expect{
				{token.String, "foobar"}, {token.Semicolon, ";"},
				{token.String, "foo bar"}, {token.Semicolon, ";"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			lexer := New(test.in)
			for _, expect := range test.expects {
				nextToken := lexer.NextToken()
				if nextToken.Type != expect.tokenType {
					t.Errorf("nextToken.Type was wrong: expected %s, but got %s\n", expect.tokenType, nextToken.Type)
				}
				if nextToken.Literal != expect.literal {
					t.Errorf("nextToken.Type was wrong: expected %s, but got %s\n", expect.literal, nextToken.Literal)
				}
			}
		})
	}
}
