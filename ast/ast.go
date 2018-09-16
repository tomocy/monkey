package ast

import (
	"fmt"
	"strings"

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
	Ident *Identifier
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
	b = append(b, s.Ident.String()...)
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
	Token token.Token
	Value Expression
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
	if s.Value != nil {
		b = append(b, s.Value.String()...)
	}
	b = append(b, ';')

	return string(b)
}

type ExpressionStatement struct {
	Token token.Token
	Value Expression
}

func (s ExpressionStatement) statement() {
}

func (s ExpressionStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s ExpressionStatement) String() string {
	if s.Value != nil {
		return s.Value.String()
	}

	return ""
}

type Integer struct {
	Token token.Token
	Value int64
}

func (i Integer) expression() {
}

func (i Integer) TokenLiteral() string {
	return i.Token.Literal
}

func (i Integer) String() string {
	return fmt.Sprint(i.Value)
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

type Boolean struct {
	Token token.Token
	Value bool
}

func (b Boolean) expression() {
}

func (b Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b Boolean) String() string {
	return fmt.Sprint(b.Value)
}

type If struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i If) expression() {
}

func (i If) TokenLiteral() string {
	return i.Token.Literal
}

func (i If) String() string {
	b := make([]byte, 0, 10)
	b = append(b, "if ("...)
	b = append(b, i.Condition.String()...)
	b = append(b, ") "...)
	b = append(b, i.Consequence.String()...)
	if i.Alternative != nil {
		b = append(b, " else "...)
		b = append(b, i.Alternative.String()...)
	}

	return string(b)
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (s BlockStatement) statement() {
}

func (s BlockStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s BlockStatement) String() string {
	b := make([]byte, 0, 10)
	b = append(b, "{ "...)
	for _, stmt := range s.Statements {
		b = append(b, stmt.String()...)
	}
	b = append(b, " }"...)

	return string(b)
}

type Function struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f Function) expression() {
}

func (f Function) TokenLiteral() string {
	return f.Token.Literal
}

func (f Function) String() string {
	b := make([]byte, 0, 10)
	b = append(b, f.TokenLiteral()...)
	b = append(b, '(')
	params := make([]string, len(f.Parameters))
	for i, param := range f.Parameters {
		params[i] = param.String()
	}
	b = append(b, strings.Join(params, ",")...)
	b = append(b, ") "...)
	b = append(b, f.Body.String()...)

	return string(b)
}

type FunctionCall struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (fc FunctionCall) expression() {
}

func (fc FunctionCall) TokenLiteral() string {
	return fc.Token.Literal
}

func (fc FunctionCall) String() string {
	b := make([]byte, 0, 10)
	b = append(b, fc.Function.String()...)
	args := make([]string, 0)
	for _, arg := range fc.Arguments {
		args = append(args, arg.String())
	}
	b = append(b, '(')
	b = append(b, strings.Join(args, ",")...)
	b = append(b, ')')

	return string(b)
}

type String struct {
	Token token.Token
	Value string
}

func (s String) expression() {
}

func (s String) TokenLiteral() string {
	return s.Token.Literal
}

func (s String) String() string {
	return s.Value
}

type Array struct {
	Token    token.Token
	Elements []Expression
}

func (a Array) expression() {
}

func (a Array) TokenLiteral() string {
	return a.Token.Literal
}

func (a Array) String() string {
	b := make([]byte, 0, 10)
	b = append(b, '[')
	elms := make([]string, len(a.Elements))
	for i, elm := range a.Elements {
		elms[i] = elm.String()
	}
	b = append(b, strings.Join(elms, ",")...)
	b = append(b, ']')

	return string(b)
}
