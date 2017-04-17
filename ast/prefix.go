package ast

import "Monkey/token"
import "bytes"

//PrefixExpression prefix
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {

}

//TokenLiteral get literal
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

//String get string
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
