package ast

import "monkey/token"

//StringLiteral "<sequence of characters>""
type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expressionNode() {

}

//TokenLiteral Pretty print
func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}

//String to string
func (s *StringLiteral) String() string {
	return s.Token.Literal
}
