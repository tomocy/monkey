package ast

import "testing"

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
		}
	}
	for _, test := range tests {
		got := Modify(test.in, turnOneIntoTwo)
		if !reflect.DeepEqual(got, test.expect) {
			t.Errorf("got was wrong: expected %+v, but got %+v\n", test.expect, got)
		}
	}
}
