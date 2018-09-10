package evaluator

import (
	"testing"

	"github.com/tomocy/monkey/lexer"
	"github.com/tomocy/monkey/object"
	"github.com/tomocy/monkey/parser"
)

func TestEvalInteger(t *testing.T) {
	tests := []struct {
		in     string
		expect int64
	}{
		{"5", 5},
		{"-5", -5},
		{"5 + 5", 10},
		{"5 - 5", 0},
		{"5 * 5", 25},
		{"5 / 5", 1},
		{"5 + 5 * 5 / 5", 10},
		{"(1 - 2) * 5 / 5", -1},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			parser := parser.New(lexer.New(test.in))
			program := parser.ParseProgram()
			got := Eval(program)
			integer, ok := got.(*object.IntegerObject)
			if !ok {
				t.Fatal("faild to assert got as *object.Integer")
			}
			if integer.Value != test.expect {
				t.Errorf("integer.Value was wrong: expected %d, but got %d\n", test.expect, integer.Value)
			}
		})
	}
}

func TestEvalBoolean(t *testing.T) {
	tests := []struct {
		in     string
		expect bool
	}{
		{"true", true},
		{"false;", false},
		{"5 < 6", true},
		{"6 < 5", false},
		{"4 < 5", true},
		{"5 < 4", false},
		{"5 < 5", false},
		{"(1 + 2) < (3 + 4)", true},
		{"(1 / 2) < (3 / 4)", false},
		{"5 == 5", true},
		{"5 != 5", false},
		{"true == true", true},
		{"true != true", false},
		{"(4 < 3) == (2 < 1)", true},
		{"(4 < 3) != (2 < 1)", false},
		{"(1 < 2) == true", true},
		{"(1 < 2) != true", false},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			parser := parser.New(lexer.New(test.in))
			program := parser.ParseProgram()
			got := Eval(program)
			boolean, ok := got.(*object.BooleanObject)
			if !ok {
				t.Fatal("faild to assert got as *object.Boolean")
			}
			if boolean.Value != test.expect {
				t.Errorf("boolean.Value was wrong: expected %t, but got %t\n", test.expect, boolean.Value)
			}
		})
	}
}

func TestBang(t *testing.T) {
	tests := []struct {
		in     string
		expect bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!!true", false},
	}
	for _, test := range tests {
		parser := parser.New(lexer.New(test.in))
		program := parser.ParseProgram()
		got := Eval(program)
		boolean, ok := got.(*object.BooleanObject)
		if !ok {
			t.Fatal("faild to assert got as *object.Boolean")
		}
		if boolean.Value != test.expect {
			t.Errorf("boolean.Value was wrong: expected %t, but got %t\n", test.expect, boolean.Value)
		}
	}
}
func TestIf(t *testing.T) {
	tests := []struct {
		in     string
		expect interface{}
	}{
		{"if (true) {10}", 10},
		{"if (false) {10}", nil},
		{"if (1 < 2) {10} else {20}", 10},
		{"if (!(1 < 2)) {10} else {20}", 20},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			parser := parser.New(lexer.New(test.in))
			program := parser.ParseProgram()
			got := Eval(program)
			if test.expect == nil {
				if got != nullObj {
					t.Errorf("got was wrong: expected %s, but got %s\n", nullObj, got)
				}
				return
			}
			expect := test.expect.(int)
			integer, ok := got.(*object.IntegerObject)
			if !ok {
				t.Fatal("faild to assert go as *object.IntegerObject")
			}
			if integer.Value != int64(expect) {
				t.Errorf("integer.Value was wrong: expected %d, but got %d\n", expect, integer.Value)
			}
		})
	}
}

func TestEvalReturnStatement(t *testing.T) {
	tests := []struct {
		in     string
		expect interface{}
	}{
		{"return 10;", 10},
		{"9; return 10;", 10},
		{"return 2*5; 11;", 10},
		{"if (true) { if (true) { return 10; } return 1; }", 10},
	}
	for _, test := range tests {
		parser := parser.New(lexer.New(test.in))
		program := parser.ParseProgram()
		got := Eval(program)
		expect := test.expect.(int)
		integer, ok := got.(*object.IntegerObject)
		if !ok {
			t.Fatal("faild to assert got as *object.IntegerObject")
		}
		if integer.Value != int64(expect) {
			t.Errorf("integer.Value was wrong: expected %d, but got %d\n", expect, integer.Value)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		in     string
		expect string
	}{
		{"if (true) { return 5 + true; }", "unknown operation: Integer + Boolean"},
		{"return true + false", "unknown operation: Boolean + Boolean"},
		{"if (!5) { return 5; }", "unknown operation: !Integer"},
		{"return -true;", "unknown operation: -Boolean"},
	}
	for _, test := range tests {
		parser := parser.New(lexer.New(test.in))
		program := parser.ParseProgram()
		got := Eval(program)
		errorObj, ok := got.(*object.ErrorObject)
		if !ok {
			t.Fatal("faild to assert got as *object.ErrorObject")
		}
		if errorObj.Message != test.expect {
			t.Errorf("errorObj.Message was wrong: expected %s, but got %s\n", test.expect, errorObj.Message)
		}
	}
}
