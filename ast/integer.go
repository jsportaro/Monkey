package ast

import "Monkey/token"

//IntegerLiteral 1,2,3,4,...
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {

}

//TokenLiteral get literal
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

//String to string
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}
