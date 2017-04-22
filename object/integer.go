package object

import (
	"fmt"
)

//Integer Integer Object
type Integer struct {
	Value int64
}

//Inspect Inspection
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

//Type type
func (i *Integer) Type() ObjectType {
	return IntegerObj
}
