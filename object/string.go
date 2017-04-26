package object

//String String Object
type String struct {
	Value string
}

//Inspect Inspection
func (s *String) Inspect() string {
	return s.Value
}

//Type type
func (s *String) Type() ObjectType {
	return StringObj
}
