package ast

import (
	"strconv"

	"github.com/tomocy/monkey/token"
)

type Node interface {
	TokenLiteral() string
	String() string
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

func (p Program) String() string {
	b := make([]byte, 0, 10)
	for _, stmt := range p.Statements {
		b = append(b, stmt.String()...)
	}

	return string(b)
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

func (s LetStatement) String() string {
	b := make([]byte, 0, 10)
	b = append(b, s.TokenLiteral()...)
	b = append(b, ' ')
	b = append(b, s.Name.String()...)
	b = append(b, " = "...)
	if s.Value != nil {
		b = append(b, s.Value.String()...)
	}
	b = append(b, ';')

	return string(b)
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

func (i Identifier) String() string {
	return i.Value
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (s ReturnStatement) statement() {
}

func (s ReturnStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s ReturnStatement) String() string {
	b := make([]byte, 0, 10)
	b = append(b, s.TokenLiteral()...)
	b = append(b, ' ')
	if s.ReturnValue != nil {
		b = append(b, s.ReturnValue.String()...)
	}
	b = append(b, ';')

	return string(b)
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (s ExpressionStatement) statement() {
}

func (s ExpressionStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s ExpressionStatement) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}

	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il IntegerLiteral) expression() {
}

func (il IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il IntegerLiteral) String() string {
	return strconv.FormatInt(il.Value, 10)
}

type Prefix struct {
	Token      token.Token
	Operator   string
	RightValue Expression
}

func (p Prefix) expression() {
}

func (p Prefix) TokenLiteral() string {
	return p.Token.Literal
}

func (p Prefix) String() string {
	b := make([]byte, 0, 10)
	b = append(b, '(')
	b = append(b, p.Operator...)
	b = append(b, p.RightValue.String()...)
	b = append(b, ')')

	return string(b)
}

type Infix struct {
	Token      token.Token
	LeftValue  Expression
	Operator   string
	RightValue Expression
}

func (i Infix) expression() {
}

func (i Infix) TokenLiteral() string {
	return i.Token.Literal
}

func (i Infix) String() string {
	b := make([]byte, 0, 10)
	b = append(b, '(')
	b = append(b, i.LeftValue.String()...)
	b = append(b, " "+i.Operator+" "...)
	b = append(b, i.RightValue.String()...)
	b = append(b, ')')

	return string(b)
}
