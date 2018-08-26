package parser

import (
	"fmt"
	"testing"

	"github.com/tomocy/monkey/ast"
	"github.com/tomocy/monkey/lexer"
	"github.com/tomocy/monkey/token"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.Let, Literal: "let"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.Ident, Literal: "foo"},
					Value: "foo",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.Ident, Literal: "bar"},
					Value: "bar",
				},
			},
		},
	}

	expected := "let foo = bar;"
	got := program.String()
	if got != expected {
		t.Errorf("the output as string from program was wrong: expected %s, but got %s\n", expected, got)
	}
}

func TestLetStatement(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	testParserHasNoErrors(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram returned nil\n")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain expected number of statements: expected %d, but got %d\n", 3, len(program.Statements))
	}

	for i, test := range tests {
		stmt := program.Statements[i]
		testLetStatement(t, stmt, test.expectedIdentifier)
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, identName string) {
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Error("faild to assert stmt as *ast.LetStatement")
	}

	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral return wrong value: expected let, but got %s\n", stmt.TokenLiteral())
	}
	if letStmt.Name.Value != identName {
		t.Errorf("letStmt.Name.Value was wrong. expected %s, but got %s\n", identName, letStmt.Name.Value)
	}
	if letStmt.Name.TokenLiteral() != identName {
		t.Errorf("letStmt.Name.TokenLiteral() returns wrong value. expected %s, but got %s\n", identName, letStmt.Name.TokenLiteral())
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 996633;
	`
	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	testParserHasNoErrors(t, parser)
	if program == nil {
		t.Fatal("parser.ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("len(program.Statement) was wrong: expected 3, but got %d\n", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Error("faild to assert stmt as *ast.ReturnStatement")
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral() returned wrong value: expected return, but got %s\n", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifier(t *testing.T) {
	input := "foobar;"
	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	testParserHasNoErrors(t, parser)
	if program == nil {
		t.Fatal("parser.ParseProgram returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) was wrong: expected 1, but got %d\n", len(program.Statements))
	}

	stmt := program.Statements[0]
	expressionStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Error("faild to assert stmt as *ast.ExpressionStatement")
	}

	ident, ok := expressionStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Error("faild to assert expressionStmt as *ast.Identifier")
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral returned wrong value: expected foobar, but got %s\n", ident.TokenLiteral())
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value was wrong: expected foobar, but got %s\n", ident.Value)
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"
	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	testParserHasNoErrors(t, parser)
	if program == nil {
		t.Fatal("program was nil")
	}
	if len(program.Statements) != 1 {
		t.Errorf("len(program.Statements) returned wrong vaue: expected 1, but got %d\n", len(program.Statements))
	}

	stmt := program.Statements[0]
	expressionStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Error("faild to assert stmt as *ast.ExpressionStatement")
	}

	intLiteral, ok := expressionStmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Error("faild to assert expressionStmt as *ast.IntegerLiteral")
	}
	if intLiteral.TokenLiteral() != "5" {
		t.Errorf("intLiteral.TokenLiteral() returned wrong value: expected 5, but got %s\n", intLiteral.TokenLiteral())
	}
	if intLiteral.Value != 5 {
		t.Errorf("intLiteral.Value was wrong: expected 5, but got %d\n", intLiteral.Value)
	}
}

func TestPrefixIntegerLiteral(t *testing.T) {
	input := `
	!5;
	-15;
	`
	tests := []struct {
		expectedOperator string
		expectedValue    int64
	}{
		{"!", 5},
		{"-", 15},
	}

	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	testParserHasNoErrors(t, parser)
	if program == nil {
		t.Fatal("program was nil")
	}
	if len(program.Statements) != 2 {
		t.Errorf("len(program.Statements) returned wrong value: expected 2, but got %d\n", len(program.Statements))
	}

	for i, test := range tests {
		stmt := program.Statements[i]
		expressionStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Error("faild to assert stmt as *ast.ExpressionStatement")
		}
		prefix, ok := expressionStmt.Expression.(*ast.Prefix)
		if !ok {
			t.Error("faild to assert expressionStmt as prefix")
		}
		if prefix.Operator != test.expectedOperator {
			t.Errorf("prefix.Operator was wrong: expected %s, but got %s\n", test.expectedOperator, prefix.Operator)
		}
		testIntegerLiteral(t, prefix.RightValue, test.expectedValue)
	}
}

func testIntegerLiteral(t *testing.T, e ast.Expression, value int64) {
	il, ok := e.(*ast.IntegerLiteral)
	if !ok {
		t.Error("faild to assert e as *ast.IntegerLiteral")
	}
	expectedTokenLiteral := fmt.Sprintf("%d", value)
	if il.TokenLiteral() != expectedTokenLiteral {
		t.Errorf("il.TokenLiteral() returned wrong value: expected %s, but got %s", expectedTokenLiteral, il.TokenLiteral())
	}
	if il.Value != value {
		t.Errorf("il.Value was wrong: expected %d, but got %d", value, il.Value)
	}
}

func testParserHasNoErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors\n", len(errs))
	for _, msg := range errs {
		t.Errorf("- %s", msg)
	}
}
