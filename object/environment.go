package object

//Environment environment
type Environment struct {
	store map[string]Object
	outer *Environment
}

//NewEnclosedEnvironment closure
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

//NewEnvironment make me the world
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

//Get recall
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

//Set store
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val

	return val
}
