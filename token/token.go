package token

const (
	_ TokenType = iota
	Illegal
	EOF

	Ident
	Int

	Assign
	Plus
	Minus
	Asterrisk
	Slash
	Bang

	Equal
	NotEqual

	LessThan
	GreaterThan

	Comma
	Semicolon

	LParen
	RParen
	LBrace
	RBrace

	Function
	Let
	If
	Else
	Return
	True
	False
)

var tokenTypes = map[string]TokenType{
	"=":    Assign,
	"+":    Plus,
	"-":    Minus,
	"*":    Asterrisk,
	"/":    Slash,
	"!":    Bang,
	"==":   Equal,
	"!=":   NotEqual,
	"<":    LessThan,
	">":    GreaterThan,
	",":    Comma,
	";":    Semicolon,
	"(":    LParen,
	")":    RParen,
	"{":    LBrace,
	"}":    RBrace,
	"\x00": EOF,
}

var keywordTypes = map[string]TokenType{
	"fn":     Function,
	"let":    Let,
	"if":     If,
	"else":   Else,
	"return": Return,
	"true":   True,
	"false":  False,
}

type TokenType int

func (tt TokenType) String() string {
	for k, v := range tokenTypes {
		if v == tt {
			return k
		}
	}

	for k, v := range keywordTypes {
		if v == tt {
			return k
		}
	}

	if tt == Illegal {
		return "Illegal"
	}

	return "Ident"
}

func LookUpTokenType(s string) TokenType {
	if tokenType, ok := tokenTypes[s]; ok {
		return tokenType
	}

	return Illegal
}

func IsKeyword(s string) bool {
	_, ok := keywordTypes[s]
	return ok
}

func LookUpKeywordType(keyword string) TokenType {
	if keywordType, ok := keywordTypes[keyword]; ok {
		return keywordType
	}

	return Illegal
}

type Token struct {
	Type    TokenType
	Literal string
}
