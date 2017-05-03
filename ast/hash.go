package ast

import (
	"monkey/token"
	"bytes"
	"strings"
)

//HashLiteral hash
type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {

}

//TokenLiteral get literal
func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

//String string
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
