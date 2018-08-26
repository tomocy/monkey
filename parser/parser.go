package parser

import (
	"github.com/tomocy/monkey/lexer"
	"github.com/tomocy/monkey/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}
