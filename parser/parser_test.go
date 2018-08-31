package parser

import (
	"testing"

	"github.com/tomocy/monkey/ast"
	"github.com/tomocy/monkey/lexer"
)

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
	type expect struct {
		tokenLiteral string
		value        int64
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"5;", expect{"5", 5}},
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

func testIntegerLiteral(t *testing.T, exp ast.Expression, expect struct {
	tokenLiteral string
	value        int64
}) {
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

func TestPrefixIntegerLiteral(t *testing.T) {
	type integerLiteral struct {
		tokenLiteral string
		value        int64
	}
	type expect struct {
		operator       string
		integerLiteral integerLiteral
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"!5;", expect{"!", integerLiteral{"5", 5}}},
		{"-5;", expect{"-", integerLiteral{"5", 5}}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testPrefix(t, expStmt.Expression, test.expect.operator)
		prefix := expStmt.Expression.(*ast.Prefix)
		testIntegerLiteral(t, prefix.RightValue, test.expect.integerLiteral)
	}
}

func testPrefix(t *testing.T, exp ast.Expression, operator string) {
	prefix, ok := exp.(*ast.Prefix)
	if !ok {
		t.Fatal("faild to assert exp as *ast.Prefix")
	}
	if prefix.Operator != operator {
		t.Errorf("prefix.Operator was wrong: expect %s, but got %s\n", operator, prefix.Operator)
	}
}

func TestInfixIntegerLiteral(t *testing.T) {
	type integerLiteral struct {
		tokenLiteral string
		value        int64
	}
	type expect struct {
		leftValue  integerLiteral
		operator   string
		rightValue integerLiteral
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"5 + 5;", expect{integerLiteral{"5", 5}, "+", integerLiteral{"5", 5}}},
		{"5 - 5;", expect{integerLiteral{"5", 5}, "-", integerLiteral{"5", 5}}},
		{"5 * 5;", expect{integerLiteral{"5", 5}, "*", integerLiteral{"5", 5}}},
		{"5 / 5;", expect{integerLiteral{"5", 5}, "/", integerLiteral{"5", 5}}},
		{"5 < 5;", expect{integerLiteral{"5", 5}, "<", integerLiteral{"5", 5}}},
		{"5 > 5;", expect{integerLiteral{"5", 5}, ">", integerLiteral{"5", 5}}},
		{"5 == 5;", expect{integerLiteral{"5", 5}, "==", integerLiteral{"5", 5}}},
		{"5 != 5;", expect{integerLiteral{"5", 5}, "!=", integerLiteral{"5", 5}}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testInfix(t, expStmt.Expression, test.expect.operator)
		infix := expStmt.Expression.(*ast.Infix)
		testIntegerLiteral(t, infix.LeftValue, test.expect.leftValue)
		testIntegerLiteral(t, infix.RightValue, test.expect.rightValue)
	}
}

func testInfix(t *testing.T, exp ast.Expression, operator string) {
	infix, ok := exp.(*ast.Infix)
	if !ok {
		t.Fatal("faild to assert exp as *ast.Infix")
	}
	if infix.Operator != operator {
		t.Errorf("infix.Operator was wrong: expected %s, but got %s", operator, infix.Operator)
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
		{"true == true;", expect{true, "==", true}},
		{"true != false;", expect{true, "!=", false}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testInfix(t, expStmt.Expression, test.expect.operator)
		infix := expStmt.Expression.(*ast.Infix)
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
	type expect struct {
		tokenLiteral string
		value        bool
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"true;", expect{"true", true}},
		{"false;", expect{"false", false}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testBoolean(t, expStmt.Expression, test.expect.value)
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
