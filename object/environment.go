package object

type Environment struct {
	objs  map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		objs: make(map[string]Object),
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

func (e Environment) Get(name string) (Object, bool) {
	obj, ok := e.objs[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) Set(name string, obj Object) {
	e.objs[name] = obj
}
