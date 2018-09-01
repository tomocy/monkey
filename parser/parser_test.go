package parser

import (
	"testing"

	"github.com/tomocy/monkey/ast"
	"github.com/tomocy/monkey/lexer"
)

type expectedInteger struct {
	tokenLiteral string
	value        int64
}

type expectedBoolean struct {
	tokenLiteral string
	value        bool
}

type expectedInfix struct {
	leftValue  expectedLiteral
	operator   string
	rightValue expectedLiteral
}

type expectedLiteral struct {
	tokenLiteral string
	value        interface{}
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
	if letStmt.Ident.TokenLiteral() != identName {
		t.Errorf("letStmt.Name.TokenLiteral() returns wrong value. expected %s, but got %s\n", identName, letStmt.Ident.TokenLiteral())
	}
	if letStmt.Ident.Value != identName {
		t.Errorf("letStmt.Name.Value was wrong. expected %s, but got %s\n", identName, letStmt.Ident.Value)
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
		testIdentifier(t, expStmt.Value, test.expect.name)
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
func TestInteger(t *testing.T) {
	tests := []struct {
		in     string
		expect expectedInteger
	}{
		{"5;", expectedInteger{"5", 5}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testInteger(t, expStmt.Value, test.expect)
	}
}
func TestPrefixInteger(t *testing.T) {
	type expect struct {
		operator   string
		rightValue expectedInteger
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"!5;", expect{"!", expectedInteger{"5", 5}}},
		{"-5;", expect{"-", expectedInteger{"5", 5}}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testPrefixInteger(t, expStmt.Value, test.expect)
	}
}

func testPrefixInteger(t *testing.T, exp ast.Expression, expect struct {
	operator   string
	rightValue expectedInteger
}) {
	prefix, ok := exp.(*ast.Prefix)
	if !ok {
		t.Fatal("faild to assert exp as *ast.Prefix")
	}
	if prefix.Operator != expect.operator {
		t.Errorf("prefix.Operator was wrong: expect %s, but got %s\n", expect.operator, prefix.Operator)
	}
	testInteger(t, prefix.RightValue, expect.rightValue)
}

func TestPrefixBoolean(t *testing.T) {
	type expect struct {
		operator   string
		rightValue expectedBoolean
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"!true", expect{"!", expectedBoolean{"true", true}}},
		{"!false", expect{"!", expectedBoolean{"false", false}}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testPrefixBoolean(t, expStmt.Value, test.expect)
	}
}

func testPrefixBoolean(t *testing.T, exp ast.Expression, expect struct {
	operator   string
	rightValue expectedBoolean
}) {
	prefix, ok := exp.(*ast.Prefix)
	if !ok {
		t.Fatal("faild to assert exp as *ast.Prefix")
	}
	if prefix.Operator != expect.operator {
		t.Errorf("prefix.Operator was wrong: expected %s, but got %s", expect.operator, prefix.Operator)
	}
	testBoolean(t, prefix.RightValue, expect.rightValue)
}
func TestInfixInteger(t *testing.T) {
	type expect struct {
		leftValue  expectedInteger
		operator   string
		rightValue expectedInteger
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"5 + 5;", expect{expectedInteger{"5", 5}, "+", expectedInteger{"5", 5}}},
		{"5 - 5;", expect{expectedInteger{"5", 5}, "-", expectedInteger{"5", 5}}},
		{"5 * 5;", expect{expectedInteger{"5", 5}, "*", expectedInteger{"5", 5}}},
		{"5 / 5;", expect{expectedInteger{"5", 5}, "/", expectedInteger{"5", 5}}},
		{"5 < 5;", expect{expectedInteger{"5", 5}, "<", expectedInteger{"5", 5}}},
		{"5 > 5;", expect{expectedInteger{"5", 5}, ">", expectedInteger{"5", 5}}},
		{"5 == 5;", expect{expectedInteger{"5", 5}, "==", expectedInteger{"5", 5}}},
		{"5 != 5;", expect{expectedInteger{"5", 5}, "!=", expectedInteger{"5", 5}}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testInfixInteger(t, expStmt.Value, test.expect)
	}
}

func testInfixInteger(t *testing.T, exp ast.Expression, expect struct {
	leftValue  expectedInteger
	operator   string
	rightValue expectedInteger
}) {
	infix, ok := exp.(*ast.Infix)
	if !ok {
		t.Fatal("faild to assert exp as *ast.Infix")
	}
	if infix.Operator != expect.operator {
		t.Errorf("infix.Operator was wrong: expected %s, but got %s", expect.operator, infix.Operator)
	}
	testInteger(t, infix.LeftValue, expect.leftValue)
	testInteger(t, infix.RightValue, expect.rightValue)
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
		testInfixBoolean(t, expStmt.Value, test.expect)
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
		testBoolean(t, expStmt.Value, test.expect)
	}
}
func TestIf(t *testing.T) {
	type expect struct {
		condition   expectedInfix
		consequence expectedLiteral
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{
			in: "if (x < y) { return x; }",
			expect: expect{
				condition:   expectedInfix{expectedLiteral{"x", "x"}, "<", expectedLiteral{"y", "y"}},
				consequence: expectedLiteral{"return", "return"},
			},
		},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		ifExp, ok := expStmt.Value.(*ast.If)
		if !ok {
			t.Fatal("faild to assert expStmt.Value as *ast.If")
		}
		testInfix(t, ifExp.Condition, test.expect.condition)
		testReturnStatement(t, ifExp.Consequence.Statements[0], test.expect.consequence.tokenLiteral)
	}
}

func TestIfElse(t *testing.T) {
	type expect struct {
		condition   expectedInfix
		consequence expectedLiteral
		alternative expectedLiteral
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{
			in: "if (x < y) { return x; } else { return y; }",
			expect: expect{
				condition:   expectedInfix{expectedLiteral{"x", "x"}, "<", expectedLiteral{"y", "y"}},
				consequence: expectedLiteral{"return", "return"},
				alternative: expectedLiteral{"return", "return"},
			},
		},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		ifExp, ok := expStmt.Value.(*ast.If)
		if !ok {
			t.Fatal("faild to assert expStmt.Value as *ast.If")
		}
		testInfix(t, ifExp.Condition, test.expect.condition)
		testReturnStatement(t, ifExp.Consequence.Statements[0], test.expect.consequence.tokenLiteral)
		testReturnStatement(t, ifExp.Alternative.Statements[0], test.expect.alternative.tokenLiteral)
	}
}

func testInfix(t *testing.T, exp ast.Expression, expect expectedInfix) {
	infix, ok := exp.(*ast.Infix)
	if !ok {
		t.Fatal("faild to assert exp as *ast.Infix")
	}
	if infix.Operator != expect.operator {
		t.Errorf("infix.Operator was wrong: expected %s, but got %s", expect.operator, infix.Operator)
	}
	testLiteral(t, infix.LeftValue, expect.leftValue)
	testLiteral(t, infix.RightValue, expect.rightValue)
}

func testLiteral(t *testing.T, exp ast.Expression, expect expectedLiteral) {
	switch v := expect.value.(type) {
	case string:
		testIdentifier(t, exp, v)
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

func testInteger(t *testing.T, exp ast.Expression, expect expectedInteger) {
	integer, ok := exp.(*ast.Integer)
	if !ok {
		t.Fatal("faild to assert exp as *ast.Integer")
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
