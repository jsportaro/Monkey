package object

//ObjectType object type
type ObjectType string

const (

	//BooleanObj bool
	BooleanObj = "BOOLEAN"
	//IntegerObj int
	IntegerObj = "INTEGER"
	//Null null
	NullObj = "NULL"
)

//Object object
type Object interface {
	Type() ObjectType
	Inspect() string
}
