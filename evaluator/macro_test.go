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
		{"quote(foo + bar)", "(foo + bar)"},
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

func TestUnquote(t *testing.T) {
	tests := []struct {
		in     string
		expect string
	}{
		{"quote(unquote(5))", "5"},
		{"quote(unquote(5 + 5))", "10"},
		{"quote(5 + unquote(2 + 3))", "(5 + 5)"},
		{"let foo = 8; quote(unquote(foo))", "8"},
		{"quote(unquote(true == false))", "false"},
		{"quote(unquote(quote(5 + 5)));", "(5 + 5)"},
		{"let quotedExp = quote(5 + 5); quote(unquote(5 + 5) + unquote(quotedExp))", "(10 + (5 + 5))"},
		{`quote(unquote("string"));`, "string"},
		{`quote(unquote([1,2,3,4]));`, "[1,2,3,4]"},
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

func TestDefineMacro(t *testing.T) {
	in := `
	let num = 1;
	let func = fn(x, y) { return x + y; };
	let m = macro(x, y) { return x + y; };
	`
	parser := parser.New(lexer.New(in))
	program := parser.ParseProgram()
	env := object.NewEnvironment()

	DefineMacros(program, env)

	if len(program.Statements) != 2 {
		t.Fatalf("len(program.Statements) returned wrong value: expected 2, but got %d\n", len(program.Statements))
	}

	if _, ok := env.Get("num"); ok {
		t.Error("invalid environment: num shuld not be defined")
	}
	if _, ok := env.Get("func"); ok {
		t.Error("invalid environment: func should not be defined")
	}

	obj, ok := env.Get("m")
	if !ok {
		t.Fatal("invalid environment: m should be defined")
	}

	macro, ok := obj.(*object.MacroObject)
	if !ok {
		t.Fatalf("assertion faild: expected *object.MacroObject, but got %T\n", obj)
	}
	if len(macro.Parameters) != 2 {
		t.Fatalf("len(macro.Parameters) returned wrong value: expected 2, but got %d\n", len(macro.Parameters))
	}
	if macro.Parameters[0].String() != "x" {
		t.Errorf("macro.Parameters[0].String() returend wrong value: expected x, but got %s\n", macro.Parameters[0].String())
	}
	if macro.Parameters[1].String() != "y" {
		t.Errorf("macro.Parameters[1].String() returend wrong value: expected y, but got %s\n", macro.Parameters[1].String())
	}
	if macro.Body.String() != "{ return (x + y); }" {
		t.Errorf("macro.Body.String() returned wront value: expected (x + y), but got %s\n", macro.Body.String())
	}
}

func TestExpandMacro(t *testing.T) {
	tests := []struct {
		in     string
		expect string
	}{
		{"let infix = macro() { quote(1 + 2); }; infix();", "(1 + 2)"},
		{"let minusReversely = macro(a, b) { quote(unquote(b) - unquote(a)); }; minusReversely(2 + 2, 10 - 5)", "((10 - 5) - (2 + 2))"},
		{
			`
			let unless = macro(condition, consequence, alternative) {
				quote(
					if (!(unquote(condition))) {
						unquote(consequence)
					} else {
						unquote(alternative)
					}
				);
			};
			unless(5 < 10, puts("not greater??"), puts("yes, greater..."))
			`,
			`if ((!(5 < 10))) { puts(not greater??) } else { puts(yes, greater...) }`,
		},
	}
	for _, test := range tests {
		parser := parser.New(lexer.New(test.in))
		program := parser.ParseProgram()
		env := object.NewEnvironment()
		DefineMacros(program, env)
		got := ExpandMacros(program, env)
		if got.String() != test.expect {
			t.Errorf("got.String() returned wrong value: expected %s, but got %s\n", test.expect, got.String())
		}
	}
}
