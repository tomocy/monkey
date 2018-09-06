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
	}
	for _, test := range tests {
		parser := parser.New(lexer.New(test.in))
		program := parser.ParseProgram()
		got := Eval(program)
		integer, ok := got.(*object.Integer)
		if !ok {
			t.Fatal("faild to assert got as *object.Integer")
		}
		if integer.Value != test.expect {
			t.Errorf("integer.Value was wrong: expected %d, but got %d\n", test.expect, integer.Value)
		}
	}
}
