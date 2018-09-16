package lexer

import (
	"github.com/tomocy/monkey/token"
)

type Lexer struct {
	input           string
	char            byte
	position        int
	readingPosition int
}

func New(input string) *Lexer {
	return &Lexer{
		input: input,
	}
}

func (l *Lexer) NextToken() token.Token {
	l.readCharacter()
	l.skipWhitespace()
	switch l.char {
	case
		'(', ')', '{', '}',
		'+', '-', '*', '/',
		'<', '>', ',', ';':
		return l.expressAsSingleToken()
	case '=', '!':
		if l.peekCharacter() == '=' {
			return l.expressAsMultipleToken()
		}

		return l.expressAsSingleToken()
	case '"':
		return l.expressAsString()
	case 0:
		return l.expressAsEOF()
	default:
		if isLetter(l.char) {
			return l.expressAsKeywordOrIdentifier()
		}
		if isDigit(l.char) {
			return l.expressAsNumber()
		}

		return l.expressAsIllegal()
	}
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.char) {
		l.readCharacter()
	}
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

func (l *Lexer) expressAsSingleToken() token.Token {
	literal := string(l.char)
	t := token.Token{
		Type:    token.LookUpTokenType(literal),
		Literal: literal,
	}

	return t
}

func (l *Lexer) expressAsMultipleToken() token.Token {
	prevChar := l.char
	l.readCharacter()
	literal := string(prevChar) + string(l.char)
	t := token.Token{
		Type:    token.LookUpTokenType(literal),
		Literal: literal,
	}

	return t
}

func (l *Lexer) expressAsString() token.Token {
	t := token.Token{
		Type:    token.String,
		Literal: l.readString(),
	}

	return t
}

func (l *Lexer) readString() string {
	beginPosition := l.readingPosition
	for {
		l.readCharacter()
		if l.char == '"' || l.char == 0 {
			break
		}
	}

	return l.input[beginPosition:l.position]
}

func (l *Lexer) expressAsEOF() token.Token {
	return token.Token{
		Type:    token.EOF,
		Literal: "",
	}
}

func (l *Lexer) expressAsKeywordOrIdentifier() token.Token {
	literal := l.readKeywordOrIdentifier()
	if token.IsKeyword(literal) {
		return token.Token{
			Type:    token.LookUpKeywordType(literal),
			Literal: literal,
		}
	}

	return token.Token{
		Type:    token.Ident,
		Literal: literal,
	}
}

func (l *Lexer) readKeywordOrIdentifier() string {
	beginPosition := l.position
	for isLetter(l.peekCharacter()) {
		l.readCharacter()
	}

	return l.input[beginPosition:l.readingPosition]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func (l *Lexer) expressAsNumber() token.Token {
	return token.Token{
		Type:    token.Int,
		Literal: l.readNumber(),
	}
}

func (l *Lexer) readNumber() string {
	beginPosition := l.position
	for isDigit(l.peekCharacter()) {
		l.readCharacter()
	}

	return l.input[beginPosition:l.readingPosition]
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (l *Lexer) expressAsIllegal() token.Token {
	return token.Token{
		Type:    token.Illegal,
		Literal: string(l.char),
	}
}

func (l *Lexer) peekCharacter() byte {
	if len(l.input) <= l.readingPosition {
		return 0
	}
	return l.input[l.readingPosition]
}

func (l *Lexer) readCharacter() {
	if len(l.input) <= l.readingPosition {
		l.char = 0
	} else {
		l.char = l.input[l.readingPosition]
	}

	l.position = l.readingPosition
	l.readingPosition++
}
