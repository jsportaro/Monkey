package object

//BuiltInFunction built in
type BuiltInFunction func(args ...Object) Object

//BuiltIn built in
type BuiltIn struct {
	Fn BuiltInFunction
}

//Type type
func (b *BuiltIn) Type() ObjectType {
	return BuiltInObj
}

//Inspect inspect
func (b *BuiltIn) Inspect() string {
	return "builtin function"
}
