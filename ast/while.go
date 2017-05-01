package ast

import "Monkey/token"
import "bytes"

//WhileExpression while (<expression>) { <statements> } will evaluate
//to the number of times the condition was evaluated
type WhileExpression struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileExpression) expressionNode() {

}

//TokenLiteral get literal
func (we *WhileExpression) TokenLiteral() string {
	return we.Token.Literal
}

//String get stringy with it
func (we *WhileExpression) String() string {
	var out bytes.Buffer
	out.WriteString("while")
	out.WriteString(" (")
	out.WriteString(we.Condition.String())
	out.WriteString(" )")
	out.WriteString(we.Body.String())

	return out.String()
}
