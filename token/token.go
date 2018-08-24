package token

const (
	Illegal = "ILLEGAL"
	EOF     = "EOF"

	Ident = "IDENT"
	Int   = "INT"

	Assign    = "="
	Plus      = "+"
	Minus     = "-"
	Asterrisk = "*"
	Slash     = "/"
	Bang      = "!"

	Eq    = "=="
	NotEq = "!="

	LT = "<"
	GT = ">"

	Comma     = ","
	Semicolon = ";"

	LParen = "("
	RParen = ")"
	LBrace = "{"
	RBrace = "}"

	Function = "FUNCTION"
	Let      = "LET"
	If       = "if"
	Else     = "else"
	Return   = "return"
	True     = "true"
	False    = "false"
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

type TokenType string

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
