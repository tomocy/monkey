package evaluator

import "github.com/tomocy/monkey/object"

type Environment struct {
	objs map[string]object.Object
}

func NewEnvironment() *Environment {
	return &Environment{
		objs: make(map[string]object.Object),
	}
}

func (e Environment) Get(name string) (object.Object, bool) {
	obj, ok := e.objs[name]
	return obj, ok
}

func (e *Environment) Set(name string, obj object.Object) {
	e.objs[name] = obj
}
