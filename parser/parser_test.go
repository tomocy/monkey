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

type expectedPrefix struct {
	operator   string
	rightValue expectedLiteral
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
		value     expectedLiteral
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"let x = 5;", expect{"x", expectedLiteral{"5", 5}}},
		{"let y = 10;", expect{"y", expectedLiteral{"10", 10}}},
		{"let foobar = foo + bar;", expect{"foobar", expectedLiteral{"foo+bar", "foo+bar"}}},
		{"let isOK = true;", expect{"isOK", expectedLiteral{"true", true}}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		testLetStatement(t, program.Statements[0], test.expect)
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, expect struct {
	identName string
	value     expectedLiteral
}) {
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Fatal("faild to assert stmt as *ast.LetStatement")
	}
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral return wrong value: expected let, but got %s\n", stmt.TokenLiteral())
	}
	if letStmt.Ident.TokenLiteral() != expect.identName {
		t.Errorf("letStmt.Name.TokenLiteral() returns wrong value. expected %s, but got %s\n", expect.identName, letStmt.Ident.TokenLiteral())
	}
	if letStmt.Ident.Value != expect.identName {
		t.Errorf("letStmt.Name.Value was wrong. expected %s, but got %s\n", expect.identName, letStmt.Ident.Value)
	}
	testLiteral(t, letStmt.Value, expect.value)
}
func TestReturnStatement(t *testing.T) {
	type expect struct {
		tokenLiteral string
		value        expectedLiteral
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{"return 5;", expect{"return", expectedLiteral{"5", 5}}},
		{"return foo;", expect{"return", expectedLiteral{"foo", "foo"}}},
		{"return true;", expect{"return", expectedLiteral{"true", true}}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		testReturnStatement(t, program.Statements[0], test.expect)
	}
}

func testReturnStatement(t *testing.T, stmt ast.Statement, expect struct {
	tokenLiteral string
	value        expectedLiteral
}) {
	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Fatal("faild to assert stmt as *ast.ReturnStatement")
	}
	if returnStmt.TokenLiteral() != expect.tokenLiteral {
		t.Errorf("returnStmt returned wrong value: expect %s, but got %s\n", expect.tokenLiteral, returnStmt.TokenLiteral())
	}
	testLiteral(t, returnStmt.Value, expect.value)
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
		t.Fatal("faild to assert exp as *ast.Identifier")
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
func TestPrefix(t *testing.T) {
	tests := []struct {
		in     string
		expect expectedPrefix
	}{
		{"!5;", expectedPrefix{"!", expectedLiteral{"5", 5}}},
		{"-5;", expectedPrefix{"-", expectedLiteral{"5", 5}}},
		{"!true", expectedPrefix{"!", expectedLiteral{"true", true}}},
		{"!false", expectedPrefix{"!", expectedLiteral{"false", false}}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testPrefix(t, expStmt.Value, test.expect)
	}
}
func TestInfix(t *testing.T) {
	tests := []struct {
		in     string
		expect expectedInfix
	}{
		{"5 + 5;", expectedInfix{expectedLiteral{"5", 5}, "+", expectedLiteral{"5", 5}}},
		{"5 - 5;", expectedInfix{expectedLiteral{"5", 5}, "-", expectedLiteral{"5", 5}}},
		{"5 * 5;", expectedInfix{expectedLiteral{"5", 5}, "*", expectedLiteral{"5", 5}}},
		{"5 / 5;", expectedInfix{expectedLiteral{"5", 5}, "/", expectedLiteral{"5", 5}}},
		{"5 < 5;", expectedInfix{expectedLiteral{"5", 5}, "<", expectedLiteral{"5", 5}}},
		{"5 > 5;", expectedInfix{expectedLiteral{"5", 5}, ">", expectedLiteral{"5", 5}}},
		{"5 == 5;", expectedInfix{expectedLiteral{"5", 5}, "==", expectedLiteral{"5", 5}}},
		{"5 != 5;", expectedInfix{expectedLiteral{"5", 5}, "!=", expectedLiteral{"5", 5}}},
		{"true == true;", expectedInfix{expectedLiteral{"true", true}, "==", expectedLiteral{"true", true}}},
		{"true != false;", expectedInfix{expectedLiteral{"true", true}, "!=", expectedLiteral{"false", false}}},
	}
	for _, test := range tests {
		parser := New(lexer.New(test.in))
		program := parser.ParseProgram()
		testParserHasNoErrors(t, parser)
		testProgramStatements(t, program.Statements, 1)
		stmt := program.Statements[0]
		testExpressionStatement(t, stmt)
		expStmt := stmt.(*ast.ExpressionStatement)
		testInfix(t, expStmt.Value, test.expect)
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
		{"1 + add(2, 3 * 4)", "(1 + add(2,(3 * 4)))"},
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
func TestIfReturn(t *testing.T) {
	type expect struct {
		condition   expectedInfix
		consequence []expectedLiteral
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{
			in: "if (x < y) { return x; }",
			expect: expect{
				condition:   expectedInfix{expectedLiteral{"x", "x"}, "<", expectedLiteral{"y", "y"}},
				consequence: []expectedLiteral{{"return", "return"}},
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
		testProgramStatements(t, ifExp.Consequence.Statements, len(test.expect.consequence))
		for i, consequenceStmt := range ifExp.Consequence.Statements {
			expectedConsequenceStmt := test.expect.consequence[i]
			testReturnStatement(t, consequenceStmt, struct {
				tokenLiteral string
				value        expectedLiteral
			}{
				"return", expectedConsequenceStmt,
			})
		}
	}
}

func TestIfElseReturn(t *testing.T) {
	type expect struct {
		condition   expectedInfix
		consequence []expectedLiteral
		alternative []expectedLiteral
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{
			in: "if (x < y) { return x; } else { return y; }",
			expect: expect{
				condition:   expectedInfix{expectedLiteral{"x", "x"}, "<", expectedLiteral{"y", "y"}},
				consequence: []expectedLiteral{{"return", "return"}},
				alternative: []expectedLiteral{{"return", "return"}},
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
		testProgramStatements(t, ifExp.Consequence.Statements, len(test.expect.consequence))
		for i, consequenceStmt := range ifExp.Consequence.Statements {
			expectedConsequenceStmt := test.expect.consequence[i]
			testReturnStatement(t, consequenceStmt, struct {
				tokenLiteral string
				value        expectedLiteral
			}{
				"return", expectedConsequenceStmt,
			})
		}
		for i, alternativeStmt := range ifExp.Consequence.Statements {
			expectedAlternativeStmt := test.expect.alternative[i]
			testReturnStatement(t, alternativeStmt, struct {
				tokenLiteral string
				value        expectedLiteral
			}{
				"return", expectedAlternativeStmt,
			})
		}
	}
}

func TestFunction(t *testing.T) {
	type expect struct {
		parameters []expectedLiteral
		body       []expectedInfix
	}
	tests := []struct {
		in     string
		expect expect
	}{
		{
			in: "fn(x, y) { x + y; }",
			expect: expect{
				parameters: []expectedLiteral{
					{"x", "x"},
					{"y", "y"},
				},
				body: []expectedInfix{
					{
						leftValue:  expectedLiteral{"x", "x"},
						operator:   "+",
						rightValue: expectedLiteral{"y", "y"},
					},
				},
			},
		},
		{
			in: "fn() { 5 + 5; }",
			expect: expect{
				parameters: make([]expectedLiteral, 0),
				body: []expectedInfix{
					{
						leftValue:  expectedLiteral{"5", 5},
						operator:   "+",
						rightValue: expectedLiteral{"5", 5},
					},
				},
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
		fn, ok := expStmt.Value.(*ast.Function)
		if !ok {
			t.Fatal("faild to assert expStmt as *ast.Function")
		}
		if len(fn.Parameters) != len(test.expect.parameters) {
			t.Errorf("len(fn.Parameters) returned wrong value: expected %d, but got %d\n", len(test.expect.parameters), len(fn.Parameters))
		}
		for i, param := range test.expect.parameters {
			testLiteral(t, fn.Parameters[i], param)
		}

		testProgramStatements(t, fn.Body.Statements, len(test.expect.body))
		for i, body := range test.expect.body {
			bodyStmt := fn.Body.Statements[i]
			testExpressionStatement(t, bodyStmt)
			bodyExpStmt := bodyStmt.(*ast.ExpressionStatement)
			testInfix(t, bodyExpStmt.Value, body)
		}
	}
}

func TestFunctionCall(t *testing.T) {
	in := "add(1, 2 * 3, 4 + 5);"
	parser := New(lexer.New(in))
	program := parser.ParseProgram()
	testParserHasNoErrors(t, parser)
	testProgramStatements(t, program.Statements, 1)
	stmt := program.Statements[0]
	testExpressionStatement(t, stmt)
	expStmt := stmt.(*ast.ExpressionStatement)
	funcCall, ok := expStmt.Value.(*ast.FunctionCall)
	if !ok {
		t.Fatal("faild to assert expStmt as *ast.FunctionCall")
	}

	testIdentifier(t, funcCall.Function, "add")

	if len(funcCall.Arguments) != 3 {
		t.Errorf("len(funcCall.Arguments) retuned wrong value: expected 3, but got %d\n", len(funcCall.Arguments))
	}
	testInteger(t, funcCall.Arguments[0], expectedInteger{"1", 1})
	testInfix(t, funcCall.Arguments[1], expectedInfix{
		expectedLiteral{"2", 2},
		"*",
		expectedLiteral{"3", 3},
	})
	testInfix(t, funcCall.Arguments[2], expectedInfix{
		expectedLiteral{"4", 4},
		"+",
		expectedLiteral{"5", 5},
	})
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
		t.Fatal("faild to assert e as *ast.Boolean")
	}
	if boolean.TokenLiteral() != expect.tokenLiteral {
		t.Errorf("boolean.TokenLiteral was wrong: expected %s, but got %s\n", expect.tokenLiteral, boolean.TokenLiteral())
	}
	if boolean.Value != expect.value {
		t.Errorf("boolean.Value was wrong: expect %t, but got %t\n", expect.value, boolean.Value)
	}
}

func testPrefix(t *testing.T, exp ast.Expression, expect expectedPrefix) {
	prefix, ok := exp.(*ast.Prefix)
	if !ok {
		t.Fatal("faild to assert exp as *ast.Prefix")
	}
	if prefix.Operator != expect.operator {
		t.Errorf("prefix.Operator was wrong: expected %s, but got %s\n", expect.operator, prefix.Operator)
	}
	testLiteral(t, prefix.RightValue, expect.rightValue)
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
	case int64:
		testInteger(t, exp, expectedInteger{expect.tokenLiteral, v})
	case bool:
		testBoolean(t, exp, expectedBoolean{expect.tokenLiteral, v})
	case string:
		testIdentifier(t, exp, v)
	}
}
