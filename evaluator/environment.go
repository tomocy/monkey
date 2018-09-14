package evaluator

import "github.com/tomocy/monkey/object"

type Environment struct {
	objs map[string]object.Object
}

func newEnvironment() *Environment {
	return &Environment{
		objs: make(map[string]object.Object),
	}
}

func (e Environment) get(name string) (object.Object, bool) {
	obj, ok := e.objs[name]
	return obj, ok
}

func (e *Environment) set(name string, obj object.Object) {
	e.objs[name] = obj
}
