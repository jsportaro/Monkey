package object

//ReturnValue return
type ReturnValue struct {
	Value Object
}

//Type type
func (rv *ReturnValue) Type() ObjectType {
	return ReturnObj
}

//Inspect inspect
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}
