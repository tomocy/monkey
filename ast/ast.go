package ast

import "github.com/tomocy/monkey/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statement()
}

type Expression interface {
	Node
	expression()
}

type Program struct {
	Statements []Statement
}

func (p Program) TokenLiteral() string {
	if 0 < len(p.Statements) {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (s LetStatement) statement() {
}

func (s LetStatement) TokenLiteral() string {
	return s.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i Identifier) expression() {
}

func (i Identifier) TokenLiteral() string {
	return i.Token.Literal
}
