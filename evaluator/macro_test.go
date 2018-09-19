package evaluator

import (
	"testing"

	"github.com/tomocy/monkey/lexer"
	"github.com/tomocy/monkey/object"
	"github.com/tomocy/monkey/parser"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		in     string
		expect string
	}{
		{"quote(5)", "5"},
		{"quote(5 + 5)", "(5 + 5)"},
		{"quote(foo)", "foo"},
		{"quote(foo + bar)", "(foo + bar)s"},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			parser := parser.New(lexer.New(test.in))
			program := parser.ParseProgram()
			env := object.NewEnvironment()
			got := Eval(program, env)
			quote, ok := got.(*object.QuoteObject)
			if !ok {
				t.Fatalf("assertion faild: expected *object.QuoteObject, but got %T\n", got)
			}
			if quote.Value == nil {
				t.Fatal("quote.Value was nil")
			}
			if quote.Value.String() != test.expect {
				t.Errorf("quote.Value.String() returned wrong value: expected %s, but got %s\n", test.expect, quote.Value.String())
			}
		})
	}
}
