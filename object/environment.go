package object

//Environment environment
type Environment struct {
	store map[string]Object
}

//NewEnvironment make me the world
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

//Get recall
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	return obj, ok
}

//Set store
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val

	return val
}
