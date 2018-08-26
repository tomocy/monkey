package parser

import (
	"fmt"

	"github.com/tomocy/monkey/ast"
	"github.com/tomocy/monkey/lexer"
	"github.com/tomocy/monkey/token"
)

type Parser struct {
	lexer          *lexer.Lexer
	currentToken   token.Token
	peekToken      token.Token
	prefixParseFns map[token.TokenType]prefixParseFunction
	infixParseFns  map[token.TokenType]infixParseFunction
	errors         []string
}

type prefixParseFunction func() ast.Expression
type infixParseFunction func(ast.Expression) ast.Expression

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:          l,
		prefixParseFns: make(map[token.TokenType]prefixParseFunction),
		infixParseFns:  make(map[token.TokenType]infixParseFunction),
		errors:         make([]string, 0),
	}
	p.nextToken()
	p.nextToken()

	return p
}

func (p Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: make([]ast.Statement, 0),
	}
	for !p.isCurrentToken(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{
		Token: p.currentToken,
	}

	if !p.isPeekToken(token.Ident) {
		p.reportPeekTokenError(token.Ident)
		return nil
	}
	p.nextToken()

	stmt.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if !p.isPeekToken(token.Assign) {
		p.reportPeekTokenError(token.Assign)
		return nil
	}

	for !p.isCurrentToken(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{
		Token: p.currentToken,
	}

	for !p.isCurrentToken(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p Parser) isPeekToken(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p Parser) isCurrentToken(tokenType token.TokenType) bool {
	return p.currentToken.Type == tokenType
}

func (p *Parser) reportPeekTokenError(tokenType token.TokenType) {
	msg := fmt.Sprintf("expected peek token to be %s, but got %s instead\n", tokenType, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}