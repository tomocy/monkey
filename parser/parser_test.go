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
	type expect struct {
		identName string
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"let x = 5;", expect{"x"}},
		{"let y = 10;", expect{"y"}},
		{"let foobar = foo + bar;", expect{"foobar"}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		testLetStatement(t, program.Statements[0], test.expect.identName)
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
	if letStmt.Name.TokenLiteral() != identName {
		t.Errorf("letStmt.Name.TokenLiteral() returns wrong value. expected %s, but got %s\n", identName, letStmt.Name.TokenLiteral())
	}
	if letStmt.Name.Value != identName {
		t.Errorf("letStmt.Name.Value was wrong. expected %s, but got %s\n", identName, letStmt.Name.Value)
	}
}

func TestReturnStatement(t *testing.T) {
	type expect struct {
		tokenLiteral string
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"return 5;", expect{"return"}},
		{"return 10;", expect{"return"}},
		{"return foo;", expect{"return"}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		testReturnStatement(t, program.Statements[0], test.expect.tokenLiteral)
	}
}

func testReturnStatement(t *testing.T, stmt ast.Statement, tokenLiteral string) {
	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Error("faild to assert stmt as *ast.ReturnStatement")
	}
	if returnStmt.TokenLiteral() != tokenLiteral {
		t.Errorf("returnStmt returned wrong value: expect %s, but got %s\n", tokenLiteral, returnStmt.TokenLiteral())
	}
}

func TestIdentifier(t *testing.T) {
	type expect struct {
		name string
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"foobar;", expect{"foobar"}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testIdentifier(t, expStmt.Expression, test.expect.name)
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, name string) {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Error("faild to assert exp as *ast.Identifier")
	}
	if ident.TokenLiteral() != name {
		t.Errorf("ident.TokenLiteral retuned wrong value: expect %s, but got %s\n", name, ident.TokenLiteral())
	}
	if ident.Value != name {
		t.Errorf("ident.Value was wrong: expect %s, but got %s\n", name, ident.TokenLiteral())
	}
}

func testProgramStatements(t *testing.T, stmts []ast.Statement, stmtLen int) {
	if len(stmts) != stmtLen {
		t.Fatalf("len(stmts) returned wrong value: expect %d, but got %d\n", stmtLen, len(stmts))
	}
}

func testExpressionStatement(t *testing.T, stmt ast.Statement) {
	if _, ok := stmt.(*ast.ExpressionStatement); !ok {
		t.Fatal("faild to assert expStmt as *ast.ExpressionStatement")
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

func TestInfixIntegerLiteral(t *testing.T) {
	input := `
	5 + 5;
	5 - 5;
	5 * 5;
	5 / 5;
	5 < 5;
	5 > 5;
	5 == 5;
	5 != 5;
	`
	tests := []struct {
		expectedLeftValue  int64
		expectedOperator   string
		expectedRightValue int64
	}{
		{5, "+", 5},
		{5, "-", 5},
		{5, "*", 5},
		{5, "/", 5},
		{5, "<", 5},
		{5, ">", 5},
		{5, "==", 5},
		{5, "!=", 5},
	}
	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	if program == nil {
		t.Fatal("program was nil")
	}
	if len(program.Statements) != len(tests) {
		t.Errorf("len(program.Statements) return wrong value: expected %d, but got %d", len(tests), len(program.Statements))
	}
	testParserHasNoErrors(t, parser)

	for i, test := range tests {
		stmt := program.Statements[i]
		expressionStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Error("faild to assert stmt as *ast.ExpressionStatement")
		}
		infix, ok := expressionStmt.Expression.(*ast.Infix)
		if !ok {
			t.Error("faild to assert expressionStmt as infix")
		}
		if infix.Operator != test.expectedOperator {
			t.Errorf("infix.Operator was wrong: expected %s, but got %s\n", test.expectedOperator, infix.Operator)
		}

		testIntegerLiteral(t, infix.LeftValue, test.expectedLeftValue)
		testIntegerLiteral(t, infix.RightValue, test.expectedRightValue)
	}
}

func TestInfixBoolean(t *testing.T) {
	type expect struct {
		leftValue  bool
		operator   string
		rightValue bool
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{
			"true == true;",
			expect{true, "==", true},
		},
		{
			"true != false;",
			expect{true, "!=", false},
		},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		if program == nil {
			t.Fatal("program was nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statements) returned wrong value: expect 1, but got %d\n", len(program.Statements))
		}
		stmt := program.Statements[0]
		expressionStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Error("falid to assert stmt as *ast.ExpressionStatement")
		}
		infix, ok := expressionStmt.Expression.(*ast.Infix)
		if !ok {
			t.Error("faild to assert expressionStmt as *ast.Infix")
		}
		if infix.Operator != test.expect.operator {
			t.Errorf("infix.Operator was wrong: expect %s, but got %s\n", test.expect.operator, infix.Operator)
		}

		testBoolean(t, infix.LeftValue, test.expect.leftValue)
		testBoolean(t, infix.RightValue, test.expect.rightValue)
	}
}

func testBoolean(t *testing.T, e ast.Expression, value bool) {
	boolean, ok := e.(*ast.Boolean)
	if !ok {
		t.Error("faild to assert e as *ast.Boolean")
	}
	if boolean.Value != value {
		t.Errorf("boolean.Value was wrong: expect %t, but got %t\n", value, boolean.Value)
	}
}

func TestPrefixAndInfixString(t *testing.T) {
	tests := []struct {
		in     string
		expect string
	}{
		{"-a * b;", "((-a) * b)"},
		{"!-a;", "(!(-a))"},
		{"a + b + c;", "((a + b) + c)"},
		{"a + b - c;", "((a + b) - c)"},
		{"a * b * c;", "((a * b) * c)"},
		{"a * b / c;", "((a * b) / c)"},
		{"a + b / c;", "(a + (b / c))"},
		{"a + b * c + d / e - f;", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5;", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4;", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4;", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 + 4 * 5;", "((3 + (4 * 5)) == (3 + (4 * 5)))"},
		{"true;", "true"},
		{"false;", "false"},
		{"3 > 5 == true;", "((3 > 5) == true)"},
		{"3 > 5 == false;", "((3 > 5) == false)"},
		{"1 + (2 + 3) + 4;", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2;", "((5 + 5) * 2)"},
		{"2 / (5 * 5);", "(2 / (5 * 5))"},
		{"-(5 + 5);", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		if program == nil {
			t.Error("program was nil")
		}
		if program.String() != test.expect {
			t.Errorf("program.String() returned wrong value: expected %s, but got %s\n", test.expect, program.String())
		}
	}
}

func TestBoolean(t *testing.T) {
	tests := []struct {
		in     string
		expect struct {
			tokenLiteral string
			value        bool
		}
	}{
		{
			"true;",
			struct {
				tokenLiteral string
				value        bool
			}{
				"true",
				true,
			},
		},
		{
			"false;",
			struct {
				tokenLiteral string
				value        bool
			}{
				"false",
				false,
			},
		},
	}

	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		if program == nil {
			t.Fatal("program was nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statements) returned wrong value: expected 1, but got %d\n", len(program.Statements))
		}
		stmt := program.Statements[0]
		expressionStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Error("faild to assert stmt as *ast.ExpressionStatement")
		}
		boolean, ok := expressionStmt.Expression.(*ast.Boolean)
		if !ok {
			t.Error("faild to assert expressionStmt as *ast.Boolean")
		}
		if boolean.TokenLiteral() != test.expect.tokenLiteral {
			t.Errorf("boolean.TokenLiteral() returned wrong value: expect %s, but got %s\n", test.expect.tokenLiteral, boolean.TokenLiteral())
		}
		if boolean.Value != test.expect.value {
			t.Errorf("boolean.Value was wrong: expect %t, but got %t\n", test.expect.value, boolean.Value)
		}
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
