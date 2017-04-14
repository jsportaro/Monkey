package ast

import (
	"Monkey/token"
	"bytes"
)

//LetStatemment let a = 0;
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

//TokenLiteral Pretty print
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) statementNode() {

}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}
