package object

import (
	"bytes"
	"strings"
)

//Array array
type Array struct {
	Elements []Object
}

//Type type
func (ao *Array) Type() ObjectType {
	return ArrayObj
}

//Inspect inspect
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}

	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
