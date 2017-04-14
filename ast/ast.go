package ast

//Node all AST elements must implement
type Node interface {
	TokenLiteral() string
	String() string
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
