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
