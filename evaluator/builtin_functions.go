package evaluator

import "github.com/tomocy/monkey/object"

var builtinFns = map[string]*object.BuiltinFunctionObject{
	"len": &object.BuiltinFunctionObject{
		Function: func(objs ...object.Object) object.Object {
			return nullObj
		},
	},
}
