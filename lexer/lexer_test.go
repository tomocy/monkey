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
				{token.Let, "let"}, {token.Ident, "five"}, {token.Assign, "="}, {token.Integer, "5"}, {token.Semicolon, ";"},
				{token.Let, "let"}, {token.Ident, "ten"}, {token.Assign, "="}, {token.Integer, "10"}, {token.Semicolon, ";"},
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
				{token.Bang, "!"}, {token.Minus, "-"}, {token.Slash, "/"}, {token.Asterrisk, "*"}, {token.Integer, "5"}, {token.Semicolon, ";"},
				{token.Integer, "5"}, {token.LessThan, "<"}, {token.Integer, "10"}, {token.GreaterThan, ">"}, {token.Integer, "5"}, {token.Semicolon, ";"},
				{token.Integer, "10"}, {token.Equal, "=="}, {token.Integer, "10"}, {token.Semicolon, ";"},
				{token.Integer, "10"}, {token.NotEqual, "!="}, {token.Integer, "10"}, {token.Semicolon, ";"},
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
				{token.If, "if"}, {token.LParen, "("}, {token.Integer, "5"}, {token.LessThan, "<"}, {token.Integer, "10"}, {token.RParen, ")"},
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
		{
			`
			$1;
			`,
			[]expect{
				{token.Illegal, "$"}, {token.Integer, "1"}, {token.Semicolon, ";"},
			},
		},
		{
			"[1, 2];",
			[]expect{
				{token.LBracket, "["}, {token.Integer, "1"}, {token.Comma, ","}, {token.Integer, "2"}, {token.RBracket, "]"}, {token.Semicolon, ";"},
			},
		},
		{
			`{"foo": "bar"};`,
			[]expect{
				{token.LBrace, "{"}, {token.String, "foo"}, {token.Colon, ":"}, {token.String, "bar"}, {token.RBrace, "}"}, {token.Semicolon, ";"},
			},
		},
		{
			"macro(x, y) { x + y; };",
			[]expect{
				{token.Macro, "macro"}, {token.LParen, "("}, {token.Ident, "x"}, {token.Comma, ","}, {token.Ident, "y"}, {token.RParen, ")"},
				{token.LBrace, "{"}, {token.Ident, "x"}, {token.Plus, "+"}, {token.Ident, "y"}, {token.Semicolon, ";"}, {token.RBrace, "}"}, {token.Semicolon, ";"},
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
