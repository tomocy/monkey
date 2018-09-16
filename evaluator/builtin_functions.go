package evaluator

import "github.com/tomocy/monkey/object"

var builtinFns = map[string]func(objs ...object.Object) object.Object{
	"len": func(objs ...object.Object) object.Object {
		return nullObj
	},
}
