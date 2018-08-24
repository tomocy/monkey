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

type Token struct {
	Type    TokenType
	Literal string
}

func LookUpIdentifier(ident string) TokenType {
	if keywordType, ok := keywordTypes[ident]; ok {
		return keywordType
	}

	return Ident
}
