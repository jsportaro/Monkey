package ast

import "Monkey/token"

//Identifier name
type Identifier struct {
	Token token.Token
	Value string
}

//TokenLiteral Pretty print
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) expressionNode() {

}

func (i *Identifier) String() string {
	return i.Value
}
