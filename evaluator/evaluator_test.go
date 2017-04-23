package evaluator

import (
	"Monkey/lexer"
	"Monkey/object"
	"Monkey/parser"
	"testing"
)

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"!true",
			false,
		},
		{
			"!false",
			true,
		},
		{
			"!5",
			false,
		},
		{
			"!!true",
			true,
		},
		{
			"!!false",
			false,
		},
		{
			"!!5",
			true,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"true",
			true,
		},
		{
			"false",
			false,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"5",
			5,
		},
		{
			"10",
			10,
		},
		{
			"-5",
			-5,
		},
		{
			"-10",
			-10,
		},
		{
			"5 + 5 + 5 + 5 - 10",
			10,
		},
		{
			"2 * 2 * 2 * 2 * 2",
			32,
		},
		{
			"(2 + 3) * 3",
			15,
		},
		{
			"2 + 3 * 3",
			11,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object isn't a Boolean.  we got %T", obj)

		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got %t but wanted %t", result.Value, expected)

		return false
	}

	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object isn't an Integer.  we got %T", obj)

		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got %d but wanted %d", result.Value, expected)

		return false
	}

	return true
}
