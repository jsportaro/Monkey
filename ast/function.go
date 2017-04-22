package ast

import (
	"Monkey/token"
	"bytes"
	"strings"
)

//FunctionLiteral fn(x) { }
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {

}

//TokenLiteral get literal
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}
