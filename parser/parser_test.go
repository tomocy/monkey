package parser

import (
	"testing"

	"github.com/tomocy/monkey/ast"
	"github.com/tomocy/monkey/lexer"
)

type expectedIntegerLiteral struct {
	tokenLiteral string
	value        int64
}

type expectedBoolean struct {
	tokenLiteral string
	value        bool
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
func TestIntegerLiteral(t *testing.T) {
	tests := []struct {
		in     string
		expect expectedIntegerLiteral
	}{
		{"5;", expectedIntegerLiteral{"5", 5}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testIntegerLiteral(t, expStmt.Expression, test.expect)
	}
}
func TestPrefixIntegerLiteral(t *testing.T) {
	type expect struct {
		operator   string
		rightValue expectedIntegerLiteral
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"!5;", expect{"!", expectedIntegerLiteral{"5", 5}}},
		{"-5;", expect{"-", expectedIntegerLiteral{"5", 5}}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testPrefixIntegerLiteral(t, expStmt.Expression, test.expect)
	}
}

func testPrefixIntegerLiteral(t *testing.T, exp ast.Expression, expect struct {
	operator   string
	rightValue expectedIntegerLiteral
}) {
	prefix, ok := exp.(*ast.Prefix)
	if !ok {
		t.Fatal("faild to assert exp as *ast.Prefix")
	}
	if prefix.Operator != expect.operator {
		t.Errorf("prefix.Operator was wrong: expect %s, but got %s\n", expect.operator, prefix.Operator)
	}
	testIntegerLiteral(t, prefix.RightValue, expect.rightValue)
}
func TestInfixIntegerLiteral(t *testing.T) {
	type expect struct {
		leftValue  expectedIntegerLiteral
		operator   string
		rightValue expectedIntegerLiteral
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"5 + 5;", expect{expectedIntegerLiteral{"5", 5}, "+", expectedIntegerLiteral{"5", 5}}},
		{"5 - 5;", expect{expectedIntegerLiteral{"5", 5}, "-", expectedIntegerLiteral{"5", 5}}},
		{"5 * 5;", expect{expectedIntegerLiteral{"5", 5}, "*", expectedIntegerLiteral{"5", 5}}},
		{"5 / 5;", expect{expectedIntegerLiteral{"5", 5}, "/", expectedIntegerLiteral{"5", 5}}},
		{"5 < 5;", expect{expectedIntegerLiteral{"5", 5}, "<", expectedIntegerLiteral{"5", 5}}},
		{"5 > 5;", expect{expectedIntegerLiteral{"5", 5}, ">", expectedIntegerLiteral{"5", 5}}},
		{"5 == 5;", expect{expectedIntegerLiteral{"5", 5}, "==", expectedIntegerLiteral{"5", 5}}},
		{"5 != 5;", expect{expectedIntegerLiteral{"5", 5}, "!=", expectedIntegerLiteral{"5", 5}}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testInfixIntegerLiteral(t, expStmt.Expression, test.expect)
	}
}

func testInfixIntegerLiteral(t *testing.T, exp ast.Expression, expect struct {
	leftValue  expectedIntegerLiteral
	operator   string
	rightValue expectedIntegerLiteral
}) {
	infix, ok := exp.(*ast.Infix)
	if !ok {
		t.Fatal("faild to assert exp as *ast.Infix")
	}
	if infix.Operator != expect.operator {
		t.Errorf("infix.Operator was wrong: expected %s, but got %s", expect.operator, infix.Operator)
	}
	testIntegerLiteral(t, infix.LeftValue, expect.leftValue)
	testIntegerLiteral(t, infix.RightValue, expect.rightValue)
}
func TestInfixBoolean(t *testing.T) {
	type expect struct {
		leftValue  expectedBoolean
		operator   string
		rightValue expectedBoolean
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"true == true;", expect{expectedBoolean{"true", true}, "==", expectedBoolean{"true", true}}},
		{"true != false;", expect{expectedBoolean{"true", true}, "!=", expectedBoolean{"false", false}}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testInfixBoolean(t, expStmt.Expression, test.expect)
	}
}

func testInfixBoolean(t *testing.T, exp ast.Expression, expect struct {
	leftValue  expectedBoolean
	operator   string
	rightValue expectedBoolean
}) {
	infix, ok := exp.(*ast.Infix)
	if !ok {
		t.Fatal("faild to assert exp as *ast.Infix")
	}
	if infix.Operator != expect.operator {
		t.Errorf("infix.Operator was wrong: expected %s, but got %s", expect.operator, infix.Operator)
	}
	testBoolean(t, infix.LeftValue, expect.leftValue)
	testBoolean(t, infix.RightValue, expect.rightValue)
}
func TestString(t *testing.T) {
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
		if program.String() != test.expect {
			t.Errorf("program.String() returned wrong value: expected %s, but got %s\n", test.expect, program.String())
		}
	}
}
func TestBoolean(t *testing.T) {
	tests := []struct {
		in     string
		expect expectedBoolean
	}{
		{"true;", expectedBoolean{"true", true}},
		{"false;", expectedBoolean{"false", false}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testBoolean(t, expStmt.Expression, test.expect)
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

func testIntegerLiteral(t *testing.T, exp ast.Expression, expect expectedIntegerLiteral) {
	integer, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Fatal("faild to assert exp as *ast.IntegerLiteral")
	}
	if integer.TokenLiteral() != expect.tokenLiteral {
		t.Errorf("integer.TokenLiteral() returned wrong value: expect %s, bot got %s\n", expect.tokenLiteral, integer.TokenLiteral())
	}
	if integer.Value != expect.value {
		t.Errorf("integer.Value was wrong: expect %d, but got %d\n", expect.value, integer.Value)
	}
}

func testBoolean(t *testing.T, e ast.Expression, expect expectedBoolean) {
	boolean, ok := e.(*ast.Boolean)
	if !ok {
		t.Error("faild to assert e as *ast.Boolean")
	}
	if boolean.TokenLiteral() != expect.tokenLiteral {
		t.Errorf("boolean.TokenLiteral was wrong: expected %s, but got %s\n", expect.tokenLiteral, boolean.TokenLiteral())
	}
	if boolean.Value != expect.value {
		t.Errorf("boolean.Value was wrong: expect %t, but got %t\n", expect.value, boolean.Value)
	}
}
