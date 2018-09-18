package object

import "testing"

func TestHashKey(t *testing.T) {
	type expect struct {
		name string
		obj1 HashKeyable
		obj2 HashKeyable
	}
	expects := []expect{
		{"Integer", &IntegerObject{Value: 5}, &IntegerObject{Value: 5}},
		{"Boolean", &BooleanObject{Value: true}, &BooleanObject{Value: true}},
		{"String", &StringObject{Value: "hello world"}, &StringObject{Value: "hello world"}},
	}
	diffObj := &StringObject{Value: "Different"}
	for _, expect := range expects {
		t.Run(expect.name, func(t *testing.T) {
			if expect.obj1.HashKey() != expect.obj2.HashKey() {
				t.Error("objects with same contents have different hash keys")
			}
			if expect.obj1.HashKey() == diffObj.HashKey() {
				t.Error("objects with different contents have same hash key")
			}
			if expect.obj2.HashKey() == diffObj.HashKey() {
				t.Error("objects with different contents have same hash key")
			}
		})
	}
}
