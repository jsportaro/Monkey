package ast

import (
	"monkey/token"
	"bytes"
	"strings"
)

//ArrayLiteral array[...]
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {

}

//TokenLiteral get literal
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

//String get string
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
