package ast

import "monkey/token"
import "bytes"

//InfixExpression infix 1 + 1
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode() {

}

//TokenLiteral get literal
func (oe *InfixExpression) TokenLiteral() string {
	return oe.Token.Literal
}

//String get string
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}
