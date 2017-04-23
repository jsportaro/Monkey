package evaluator

import (
	"Monkey/ast"
	"Monkey/object"
)

var (
	//NULL null
	NULL = &object.Null{}
	//TRUE true
	TRUE = &object.Boolean{Value: true}
	//FALSE false
	FALSE = &object.Boolean{Value: false}
)

//Eval eval ast
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	}

	return nil
}

func evaluateStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}
