package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() Expression { return &Integer{Value: 1} }
	two := func() Expression { return &Integer{Value: 2} }
	turnOneIntoTwo := func(node Node) Node {
		integer, ok := node.(*Integer)
		if !ok {
			return node
		}
		if integer.Value != 1 {
			return node
		}
		integer.Value = 2
		return integer
	}
	tests := []struct {
		in     Node
		expect Node
	}{
		{one(), two()},
		{
			&Program{
				[]Statement{&ExpressionStatement{Value: one()}},
			},
			&Program{
				[]Statement{&ExpressionStatement{Value: two()}},
			},
		},
		{
			&LetStatement{Ident: &Identifier{Value: "x"}, Value: one()},
			&LetStatement{Ident: &Identifier{Value: "x"}, Value: two()},
		},
		{
			&ReturnStatement{Value: one()},
			&ReturnStatement{Value: two()},
		},
		{
			&If{
				Condition: one(),
				Consequence: &BlockStatement{
					Statements: []Statement{&ExpressionStatement{Value: one()}},
				},
				Alternative: &BlockStatement{
					Statements: []Statement{&ExpressionStatement{Value: one()}},
				},
			},
			&If{
				Condition: two(),
				Consequence: &BlockStatement{
					Statements: []Statement{&ExpressionStatement{Value: two()}},
				},
				Alternative: &BlockStatement{
					Statements: []Statement{&ExpressionStatement{Value: two()}},
				},
			},
		},
		{
			&Prefix{Operator: "-", RightValue: one()},
			&Prefix{Operator: "-", RightValue: two()},
		},
		{
			&Infix{LeftValue: one(), Operator: "+", RightValue: one()},
			&Infix{LeftValue: two(), Operator: "+", RightValue: two()},
		},
		{
			&Subscript{LeftValue: one(), Index: one()},
			&Subscript{LeftValue: two(), Index: two()},
		},
	}
	for _, test := range tests {
		t.Run(test.in.String(), func(t *testing.T) {
			got := Modify(test.in, turnOneIntoTwo)
			if !reflect.DeepEqual(got, test.expect) {
				t.Errorf("got was wrong: expected %+v, but got %+v\n", test.expect, got)
			}
		})
	}
}
