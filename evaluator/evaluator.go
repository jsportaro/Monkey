package evaluator

import (
	"Monkey/ast"
	"Monkey/object"
	"fmt"
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
		return evaluateProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evaluatePrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evaluateInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evaluateBlockStatement(node)
	case *ast.IfExpression:
		return evaluateIfExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
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

func evaluateBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil {
			rt := result.Type()
			if rt == object.ReturnObj || rt == object.ErrorObj {
				return result
			}
		}
	}

	return result
}

func evaluateIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isError(condition) {
		return condition
	}
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
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evaluateNegationOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerObj {
		return newError("unknown operator: -%s", right.Type())
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
		return newError("unknown oeprator: %s%s", operator, right.Type())
	}
}

func evaluateProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result

}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObj
	}

	return false
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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
