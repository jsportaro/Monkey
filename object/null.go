package object

//Null null
type Null struct {
}

//Inspect inspect
func (n *Null) Inspect() string {
	return "null"
}

//Type type
func (n *Null) Type() ObjectType {
	return NullObj
}
