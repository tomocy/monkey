package parser

import (
	"testing"

	"github.com/tomocy/monkey/ast"
	"github.com/tomocy/monkey/lexer"
)

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

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
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
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral return wrong value: expected let, but got %s\n", stmt.TokenLiteral())
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Error("faild to assert stmt as *ast.LetStatement\n")
	}

	if letStmt.Name.Value != identName {
		t.Errorf("letStmt.Name.Value was wrong. expected %s, but got %s\n", identName, letStmt.Name.Value)
	}
	if letStmt.Name.TokenLiteral() != identName {
		t.Errorf("letStmt.Name.TokenLiteral() returns wrong value. expected %s, but got %s\n", identName, letStmt.Name.TokenLiteral())
	}
}
