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
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evaluatePrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evaluateInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evaluateStatements(node.Statements)
	case *ast.IfExpression:
		return evaluateIfExpression(node)
	}

	return nil
}

func evaluateBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evaluateIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func evaluateInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
		return evaluateIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL
	}

}

func evaluateIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftInt := left.(*object.Integer).Value
	rightInt := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftInt + rightInt}
	case "-":
		return &object.Integer{Value: leftInt - rightInt}
	case "*":
		return &object.Integer{Value: leftInt * rightInt}
	case "/":
		return &object.Integer{Value: leftInt / rightInt}
	case "<":
		return nativeBoolToBooleanObject(leftInt < rightInt)
	case ">":
		return nativeBoolToBooleanObject(leftInt > rightInt)
	case "==":
		return nativeBoolToBooleanObject(leftInt == rightInt)
	case "!=":
		return nativeBoolToBooleanObject(leftInt != rightInt)
	default:
		return NULL
	}
}

func evaluateNegationOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerObj {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evaluatePrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evaluateBangOperatorExpression(right)
	case "-":
		return evaluateNegationOperatorExpression(right)
	default:
		return NULL
	}
}

func evaluateStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
	}

	return result

}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}
