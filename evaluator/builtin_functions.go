package evaluator

import "github.com/tomocy/monkey/object"

var builtinFns = map[string]*object.BuiltinFunctionObject{
	"len": &object.BuiltinFunctionObject{
		Function: func(objs ...object.Object) object.Object {
			if len(objs) != 1 {
				return newError("too many arguments to len: expected 1, but got %d", len(objs))
			}
			switch obj := objs[0].(type) {
			case *object.StringObject:
				return &object.IntegerObject{Value: int64(len(obj.Value))}
			case *object.ArrayObject:
				return &object.IntegerObject{Value: int64(len(obj.Elements))}
			default:
				return newError("unknown operation: len(%s)", obj.Type())
			}
		},
	},
	"first": &object.BuiltinFunctionObject{
		Function: func(objs ...object.Object) object.Object {
			if len(objs) != 1 {
				return newError("invalid number of arguments to first: expected 1, but got %d", len(objs))
			}

			obj := objs[0]
			array, ok := obj.(*object.ArrayObject)
			if !ok {
				return newError("unknown operation: first(%s)", obj.Type())
			}

			if len(array.Elements) <= 0 {
				return nullObj
			}

			return array.Elements[0]
		},
	},
}
