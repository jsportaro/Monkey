package ast

import "Monkey/token"
import "bytes"

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {

}

//TokenLiteral get literal
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

//String get stringy with it
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
