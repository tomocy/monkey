package object

import "testing"

func TestHashKey(t *testing.T) {
	type expect struct {
		obj1 Object
		obj2 Object
	}
	expects := []expect{
		{&IntegerObject{Value: 5}, &IntegerObject{Value: 1}},
		{&BooleanObject{Value: true}, &BooleanObject{Value: true}},
		{&StringObject{Value: "hello world"}, &StringObject{Value: "hello world"}},
	}
	diffObj := &StringObject{Value: "Different"}
	for _, expect := range expects {
		t.Run(expect.obj1.Inspect(), func(t *testing.T) {
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
