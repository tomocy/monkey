package object

var builtinFns = map[string]builtinFunction{
	"len": func(objs ...Object) Object {
		return &NullObject{}
	},
}
