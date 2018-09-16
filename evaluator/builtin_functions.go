package evaluator

import "github.com/tomocy/monkey/object"

var builtinFns = map[string]*object.BuiltinFunctionObject{
	"len": &object.BuiltinFunctionObject{
		Function: func(objs ...object.Object) object.Object {
			if len(objs) != 1 {
				return newError("too many arguments to len: expected 1, but got %d", len(objs))
			}
			obj := objs[0]
			str, ok := obj.(*object.StringObject)
			if !ok {
				return newError("unknown operation: len(%s)", obj.Type())
			}

			return &object.IntegerObject{Value: int64(len(str.Value))}
		},
	},
}
