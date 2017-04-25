package object

//ObjectType object type
type ObjectType string

const (

	//BooleanObj bool
	BooleanObj = "BOOLEAN"
	//IntegerObj int
	IntegerObj = "INTEGER"
	//NullObj null
	NullObj = "NULL"
	//ReturnObj return
	ReturnObj = "RETURN"
	//ErrorObj error
	ErrorObj = "ERROR"
	//FunctionObj function
	FunctionObj = "FUNCTION"
)

//Object object
type Object interface {
	Type() ObjectType
	Inspect() string
}
