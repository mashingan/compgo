package interp

type Environment map[string]Object

func NewEnvironment() Environment {
	return Environment{}
}

func (e Environment) Get(name string) (Object, bool) {
	o, ok := e[name]
	return o, ok
}

func (e Environment) Set(name string, val Object) Object {
	e[name] = val
	return val
}
