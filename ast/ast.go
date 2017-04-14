package ast

import "Monkey/token"

//Node all AST elements must implement
type Node interface {
	TokenLiteral() string
}

//Statement Doesn't produce a valut
type Statement interface {
	Node
	statementNode()
}

//Expression Produce a value
type Expression interface {
	Node
	expressionNode()
}

//Program Collection of statements
type Program struct {
	Statements []Statement
}

//TokenLiteral Gets the string representation
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

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
