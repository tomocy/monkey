package object

type Environment struct {
	objs map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{
		objs: make(map[string]Object),
	}
}

func (e Environment) Get(name string) (Object, bool) {
	obj, ok := e.objs[name]
	return obj, ok
}

func (e *Environment) Set(name string, obj Object) {
	e.objs[name] = obj
}
