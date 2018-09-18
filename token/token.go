package token

const (
	Illegal = "Illegal"
	EOF     = "EOF"

	Ident = "Ident"

	Assign    = "Assign"
	Plus      = "Plus"
	Minus     = "Minus"
	Asterrisk = "Asterrisk"
	Slash     = "Slash"
	Bang      = "Bang"

	Equal    = "Equal"
	NotEqual = "NotEqual"

	LessThan    = "LessThan"
	GreaterThan = "GreaterThan"

	Comma     = "Comma"
	Colon     = "Colon"
	Semicolon = "Semicolon"

	LParen   = "LParen"
	RParen   = "RParen"
	LBrace   = "LBrace"
	RBrace   = "RBrace"
	LBracket = "LBracket"
	RBracket = "RBracket"

	Function = "Function"
	Let      = "Let"
	If       = "If"
	Else     = "Else"
	Return   = "Return"
	True     = "True"
	False    = "False"
	Int      = "Int"
	String   = "String"
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
	":":    Colon,
	";":    Semicolon,
	"(":    LParen,
	")":    RParen,
	"{":    LBrace,
	"}":    RBrace,
	"[":    LBracket,
	"]":    RBracket,
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

type TokenType string

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
