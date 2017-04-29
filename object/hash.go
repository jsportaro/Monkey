package object

import (
	"bytes"
	"fmt"
	"strings"
)

//HashPair hash pair
type HashPair struct {
	Key   Object
	Value Object
}

//Hash hash
type Hash struct {
	Pairs map[HashKey]HashPair
}

//Type type
func (h *Hash) Type() ObjectType { return HashObj }

//Inspect inspect
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
