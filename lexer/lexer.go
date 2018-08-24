package lexer

import "github.com/tomocy/monkey/token"

type Lexer struct {
	input           string
	char            byte
	position        int
	readingPosition int
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()

	return l
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token
	l.skipWhitespace()
	switch l.char {
	case '=':
		if l.peekChar() == '=' {
			prevChar := l.char
			l.readChar()
			t.Type = token.EQ
			t.Literal = string(prevChar) + string(l.char)
		} else {
			t = newToken(token.ASSIGN, l.char)
		}
	case '+':
		t = newToken(token.PLUS, l.char)
	case '-':
		t = newToken(token.MINUS, l.char)
	case '*':
		t = newToken(token.ASTERRISK, l.char)
	case '/':
		t = newToken(token.SLASH, l.char)
	case '!':
		if l.peekChar() == '=' {
			prevChar := l.char
			l.readChar()
			t.Type = token.NOT_EQ
			t.Literal = string(prevChar) + string(l.char)
		} else {
			t = newToken(token.BANG, l.char)
		}
	case '<':
		t = newToken(token.LT, l.char)
	case '>':
		t = newToken(token.GT, l.char)
	case ',':
		t = newToken(token.COMMA, l.char)
	case ';':
		t = newToken(token.SEMICOLON, l.char)
	case '(':
		t = newToken(token.LPAREN, l.char)
	case ')':
		t = newToken(token.RPAREN, l.char)
	case '{':
		t = newToken(token.LBRACE, l.char)
	case '}':
		t = newToken(token.RBRACE, l.char)
	case 0:
		t.Type = token.EOF
		t.Literal = ""
	default:
		if isLetter(l.char) {
			t.Literal = l.readIdentifier()
			t.Type = token.LookUpIdentifier(t.Literal)
			return t
		}

		if isDigit(l.char) {
			t.Literal = l.readNumber()
			t.Type = token.INT
			return t
		}

		t = newToken(token.ILLEGAL, l.char)
	}

	l.readChar()
	return t
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.char) {
		l.readChar()
	}
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

func (l *Lexer) peekChar() byte {
	if len(l.input) <= l.readingPosition {
		return 0
	}
	return l.input[l.readingPosition]
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(char),
	}
}

func (l *Lexer) readIdentifier() string {
	beginPosition := l.position
	for isLetter(l.char) {
		l.readChar()
	}

	return l.input[beginPosition:l.position]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func (l *Lexer) readNumber() string {
	beginPosition := l.position
	for isDigit(l.char) {
		l.readChar()
	}

	return l.input[beginPosition:l.position]
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (l *Lexer) readChar() {
	if len(l.input) <= l.readingPosition {
		l.char = 0
	} else {
		l.char = l.input[l.readingPosition]
	}

	l.position = l.readingPosition
	l.readingPosition++
}
