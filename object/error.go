package object

//Error error
type Error struct {
	Message string
}

//Type type
func (e *Error) Type() ObjectType {
	return ErrorObj
}

//Inspect inspect
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}
