package interp

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{map[string]Object{}, nil}
}

func NewEnvironmentFrame(parent *Environment) *Environment {
	return &Environment{map[string]Object{}, parent}
}

func (e Environment) Get(name string) (Object, bool) {
	o, ok := e.store[name]
	if !ok && e.outer != nil {
		oo, ook := e.outer.Get(name)
		return oo, ook
	}
	return o, ok
}

func (e Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
