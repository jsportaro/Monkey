package ast

import "monkey/token"

//Boolean true/false
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {
}

//TokenLiteral get literal
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

//String get stringy with it
func (b *Boolean) String() string {
	return b.TokenLiteral()
}
