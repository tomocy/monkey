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
	Subscript
)

var precedences = map[token.TokenType]precedence{
	token.Equal:       Equal,
	token.NotEqual:    Equal,
	token.LessThan:    Relational,
	token.GreaterThan: Relational,
	token.Plus:        Additive,
	token.Minus:       Additive,
	token.Asterrisk:   Multiplicative,
	token.Slash:       Multiplicative,
	token.LParen:      Call,
	token.LBracket:    Subscript,
}

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
	p.registerPrefixParseFunction(token.Integer, p.parseInterger)
	p.registerPrefixParseFunction(token.Bang, p.parsePrefix)
	p.registerPrefixParseFunction(token.Minus, p.parsePrefix)
	p.registerPrefixParseFunction(token.True, p.parseBoolean)
	p.registerPrefixParseFunction(token.False, p.parseBoolean)
	p.registerPrefixParseFunction(token.LParen, p.parseGroupedExpression)
	p.registerPrefixParseFunction(token.If, p.parseIf)
	p.registerPrefixParseFunction(token.Function, p.parseFunction)
	p.registerPrefixParseFunction(token.String, p.parseString)
	p.registerPrefixParseFunction(token.LBracket, p.parseArray)
	p.registerPrefixParseFunction(token.LBrace, p.parseHash)
	p.registerPrefixParseFunction(token.Macro, p.parseMacro)

	p.registerInfixParseFunction(token.Equal, p.parseInfix)
	p.registerInfixParseFunction(token.NotEqual, p.parseInfix)
	p.registerInfixParseFunction(token.LessThan, p.parseInfix)
	p.registerInfixParseFunction(token.GreaterThan, p.parseInfix)
	p.registerInfixParseFunction(token.Plus, p.parseInfix)
	p.registerInfixParseFunction(token.Minus, p.parseInfix)
	p.registerInfixParseFunction(token.Asterrisk, p.parseInfix)
	p.registerInfixParseFunction(token.Slash, p.parseInfix)
	p.registerInfixParseFunction(token.LParen, p.parseFunctionCall)
	p.registerInfixParseFunction(token.LBracket, p.parseSubscript)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefixParseFunction(tokenType token.TokenType, fn prefixParseFunction) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixParseFunction(tokenType token.TokenType, fn infixParseFunction) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseInterger() ast.Expression {
	value, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %s as int64\n", p.currentToken.Literal))
		return nil
	}

	return &ast.Integer{
		Token: p.currentToken,
		Value: value,
	}
}

func (p *Parser) parsePrefix() ast.Expression {
	exp := &ast.Prefix{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}
	p.nextToken()
	exp.RightValue = p.parseExpression(Prefix)

	return exp
}

func (p *Parser) parseInfix(leftValue ast.Expression) ast.Expression {
	exp := &ast.Infix{
		Token:     p.currentToken,
		LeftValue: leftValue,
		Operator:  p.currentToken.Literal,
	}
	prec := p.currentPrecedence()

	p.nextToken()
	exp.RightValue = p.parseExpression(prec)

	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.isCurrentToken(token.True),
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(Lowest)
	if !p.isPeekToken(token.RParen) {
		return nil
	}
	p.nextToken()

	return exp
}

func (p *Parser) parseIf() ast.Expression {
	exp := &ast.If{
		Token: p.currentToken,
	}
	if !p.isPeekToken(token.LParen) {
		p.reportPeekTokenError(token.LParen)
		return nil
	}
	p.nextToken()
	p.nextToken()

	exp.Condition = p.parseExpression(Lowest)
	if !p.isPeekToken(token.RParen) {
		p.reportPeekTokenError(token.RParen)
		return nil
	}
	p.nextToken()

	if !p.isPeekToken(token.LBrace) {
		p.reportPeekTokenError(token.LBrace)
		return nil
	}
	p.nextToken()

	exp.Consequence = p.parseBlockStatement()

	if p.isPeekToken(token.Else) {
		p.nextToken()
		if !p.isPeekToken(token.LBrace) {
			p.reportPeekTokenError(token.LBrace)
			return nil
		}
		p.nextToken()
		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseFunction() ast.Expression {
	exp := &ast.Function{
		Token: p.currentToken,
	}

	if !p.isPeekToken(token.LParen) {
		p.reportPeekTokenError(token.LParen)
		return nil
	}
	p.nextToken()
	exp.Parameters = p.parseFunctionParameters()

	if !p.isPeekToken(token.LBrace) {
		p.reportPeekTokenError(token.LBrace)
		return nil
	}
	p.nextToken()
	exp.Body = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseFunctionCall(function ast.Expression) ast.Expression {
	exp := &ast.FunctionCall{
		Token:    p.currentToken,
		Function: function,
	}
	exp.Arguments = p.parseFunctionCallArguments()

	return exp
}

func (p *Parser) parseFunctionCallArguments() []ast.Expression {
	p.nextToken()
	args := make([]ast.Expression, 0)
	if p.isCurrentToken(token.RParen) {
		return args
	}

	args = append(args, p.parseExpression(Lowest))
	for p.isPeekToken(token.Comma) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(Lowest))
	}

	if !p.isPeekToken(token.RParen) {
		p.reportPeekTokenError(token.RParen)
		return nil
	}

	p.nextToken()

	return args
}

func (p *Parser) parseString() ast.Expression {
	return &ast.String{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseArray() ast.Expression {
	exp := &ast.Array{Token: p.currentToken}
	exp.Elements = p.parseArrayElements()

	return exp
}

func (p *Parser) parseArrayElements() []ast.Expression {
	p.nextToken()
	exps := make([]ast.Expression, 0)
	if p.isCurrentToken(token.RBracket) {
		return exps
	}
	// [1, 2]
	// [1, 2]
	exps = append(exps, p.parseExpression(Lowest))
	for p.isPeekToken(token.Comma) {
		p.nextToken()
		p.nextToken()
		exps = append(exps, p.parseExpression(Lowest))
	}

	if !p.isPeekToken(token.RBracket) {
		p.reportPeekTokenError(token.RBracket)
		return nil
	}

	p.nextToken()

	return exps
}

func (p *Parser) parseSubscript(leftValue ast.Expression) ast.Expression {
	exp := &ast.Subscript{
		Token:     p.currentToken,
		LeftValue: leftValue,
	}
	p.nextToken()

	exp.Index = p.parseExpression(Lowest)

	if !p.isPeekToken(token.RBracket) {
		p.reportPeekTokenError(token.RBracket)
		return nil
	}

	p.nextToken()

	return exp
}

func (p *Parser) parseHash() ast.Expression {
	exp := &ast.Hash{
		Token:  p.currentToken,
		Values: make(map[ast.Expression]ast.Expression),
	}

	p.nextToken()
	if p.isCurrentToken(token.RBrace) {
		return exp
	}

	key := p.parseExpression(Lowest)
	if !p.isPeekToken(token.Colon) {
		p.reportPeekTokenError(token.Colon)
		return nil
	}
	p.nextToken()
	p.nextToken()
	value := p.parseExpression(Lowest)
	exp.Values[key] = value

	for p.isPeekToken(token.Comma) {
		p.nextToken()
		p.nextToken()
		key := p.parseExpression(Lowest)
		if !p.isPeekToken(token.Colon) {
			p.reportPeekTokenError(token.Colon)
			return nil
		}
		p.nextToken()
		p.nextToken()
		value := p.parseExpression(Lowest)
		exp.Values[key] = value
	}

	if !p.isPeekToken(token.RBrace) {
		p.reportPeekTokenError(token.RBrace)
		return nil
	}

	p.nextToken()

	return exp
}

func (p *Parser) parseMacro() ast.Expression {
	exp := &ast.Macro{
		Token: p.currentToken,
	}

	if !p.isPeekToken(token.LParen) {
		p.reportPeekTokenError(token.LParen)
		return nil
	}
	p.nextToken()
	exp.Parameters = p.parseFunctionParameters()

	if !p.isPeekToken(token.LBrace) {
		p.reportPeekTokenError(token.LBrace)
		return nil
	}
	p.nextToken()
	exp.Body = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	p.nextToken()
	idents := make([]*ast.Identifier, 0)
	if p.isCurrentToken(token.RParen) {
		return idents
	}

	ident := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
	idents = append(idents, ident)
	for p.isPeekToken(token.Comma) {
		p.nextToken()
		p.nextToken()
		ident = &ast.Identifier{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		}
		idents = append(idents, ident)
	}

	if !p.isPeekToken(token.RParen) {
		p.reportPeekTokenError(token.RParen)
		return nil
	}

	p.nextToken()

	return idents
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStmt := &ast.BlockStatement{
		Token:      p.currentToken,
		Statements: make([]ast.Statement, 0),
	}
	p.nextToken()

	for !p.isCurrentToken(token.RBrace) && !p.isCurrentToken(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			blockStmt.Statements = append(blockStmt.Statements, stmt)
		}
		p.nextToken()
	}

	return blockStmt
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

	stmt.Ident = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if !p.isPeekToken(token.Assign) {
		p.reportPeekTokenError(token.Assign)
		return nil
	}

	p.nextToken()
	p.nextToken()

	stmt.Value = p.parseExpression(Lowest)

	if p.isPeekToken(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{
		Token: p.currentToken,
	}

	p.nextToken()

	stmt.Value = p.parseExpression(Lowest)

	if p.isPeekToken(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{
		Token: p.currentToken,
		Value: p.parseExpression(Lowest),
	}

	if p.isPeekToken(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence precedence) ast.Expression {
	prefixParseFn := p.prefixParseFns[p.currentToken.Type]
	if prefixParseFn == nil {
		p.reportNoPrefixParseFunction(p.currentToken.Type)
		return nil
	}
	leftValue := prefixParseFn()

	for !p.isPeekToken(token.Semicolon) && precedence < p.peekPrecedence() {
		infixParseFn := p.infixParseFns[p.peekToken.Type]
		if infixParseFn == nil {
			p.reportNoInfixParseFunction(p.peekToken.Type)
			return leftValue
		}

		p.nextToken()
		leftValue = infixParseFn(leftValue)
	}

	return leftValue
}

func (p Parser) isPeekToken(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p Parser) isCurrentToken(tokenType token.TokenType) bool {
	return p.currentToken.Type == tokenType
}

func (p Parser) peekPrecedence() precedence {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return Lowest
}

func (p Parser) currentPrecedence() precedence {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return Lowest
}

func (p *Parser) reportPeekTokenError(tokenType token.TokenType) {
	msg := fmt.Sprintf("expected peek token to be %s, but got %s instead\n", tokenType, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) reportNoPrefixParseFunction(tokenType token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", tokenType)
	p.errors = append(p.errors, msg)
}

func (p *Parser) reportNoInfixParseFunction(tokenType token.TokenType) {
	msg := fmt.Sprintf("no infix parse function for %s found", tokenType)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p Parser) Errors() []string {
	return p.errors
}
