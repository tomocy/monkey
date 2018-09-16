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
}
