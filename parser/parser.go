package parser

import (
	"fmt"
	"strconv"

	"github.com/tomocy/monkey/ast"
	"github.com/tomocy/monkey/lexer"
	"github.com/tomocy/monkey/token"
)

const (
	_ precedence = iota
	Lowest
	Equal
	Relational
	Additive
	Multiplicative
	Prefix
	Call
)

type precedence int

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

	p.registerPrefixParseFunction(token.Ident, p.parseIdentifier)
	p.registerPrefixParseFunction(token.Int, p.parseIntergerLiteral)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefixParseFunction(tokenType token.TokenType, fn prefixParseFunction) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseIntergerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %s as int64\n", p.currentToken.Literal))
		return nil
	}

	return &ast.IntegerLiteral{
		Token: p.currentToken,
		Value: value,
	}
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
		return p.parseExpressionStatement()
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

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{
		Token:      p.currentToken,
		Expression: p.parseExpression(Lowest),
	}

	if p.isPeekToken(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence precedence) ast.Expression {
	prefixParseFn := p.prefixParseFns[p.currentToken.Type]
	if prefixParseFn == nil {
		return nil
	}

	return prefixParseFn()
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

func (p Parser) Errors() []string {
	return p.errors
}
