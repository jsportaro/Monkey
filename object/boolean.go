package object

import (
	"fmt"
)

//Boolean boolean
type Boolean struct {
	Value bool
}

//Inspect inspect
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

//Type type
func (b *Boolean) Type() ObjectType {
	return BooleanObj
}
