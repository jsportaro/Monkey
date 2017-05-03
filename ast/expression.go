package ast

import "monkey/token"

//ExpressionStatement Example: x + 5;
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {

}

//TokenLiteral Get literal
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

//String Get String
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}
