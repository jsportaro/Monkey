package evaluator

import (
	"Monkey/lexer"
	"Monkey/object"
	"Monkey/parser"
	"testing"
)

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i =0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1]",
			3,
		},
		{
			"let a = [1, 2, 3]; a[2]",
			3,
		},
		{
			"let a = [1, 2, 3]; a[0] + a[1] + a[2];",
			6,
		},
		{
			"let a = [1, 2, 3]; let i = a[0]; a[i];",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			if !testIntegerObject(t, evaluated, int64(integer)) {
				t.Errorf("For %s", tt.input)
			}
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)

	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("Expected *object.Array but got %T", evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong number of elements, got %d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

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

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`len([1, 2, 3])`, 3},
		{`first([1, 2, 3])`, 1},
		{`first()`, "wrong number of arguments. wanted 1 got 0"},
		{`first("abc")`, "arguments to `first` must be ARRAY"},
		{`last([1, 2, 3])`, 3},
		{`last()`, "wrong number of arguments. wanted 1 got 0"},
		{`last("abc")`, "arguments to `last` must be ARRAY"},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest()`, "wrong number of arguments. wanted 1 got 0"},
		{`rest("abc")`, "arguments to `rest` must be ARRAY"},
		{`push([1, 2, 3], 1)`, []int{1, 2, 3, 1}},
		{`push()`, "wrong number of arguments. wanted 2 got 0"},
		{`push("abc", 1)`, "first argument to `push` must be ARRAY"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {

		case int:
			if !testIntegerObject(t, evaluated, int64(expected)) {
				t.Errorf("Error for test %s wanted %d got %s", tt.input, expected, evaluated)
			}
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("for %q object is not Error, got %T", tt.input, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testIntegerObject(t, array.Elements[i], int64(expectedElem))
			}
		}
	}
}

func TestClosure(t *testing.T) {
	input := `
	let newAdder = fn(x) {
			fn(y) { x + y; }
	}
	
	let addTwo = newAdder(2);
	
	addTwo(2);`

	testIntegerObject(t, testEval(input), 4)
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}

				return 1;
			}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errorObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got %T", evaluated)
			continue
		}

		if errorObj.Message != tt.expected {
			t.Errorf("Wong error message.  wanted %q but got %q", tt.expected, errorObj.Message)
		}

	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"let identity = fn(x) { x; }; identity(5);",
			5,
		},
		{
			"let identity = fn(x) { return x; }; identity(5);",
			5,
		},
		{
			"let double = fn(x) { x * 2; }; double(5);",
			10,
		},
		{
			"let add = fn(x, y) { x + y; }; add(5, 5);",
			10,
		},
		{
			"let add = fn(x, y) { x + y; }; add(5 + 5, add(5,5));",
			20,
		},
		{
			"fn(x){ x; }(5)",
			5,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
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
		{
			"1 < 2",
			true,
		},
		{
			"1 > 2",
			false,
		},
		{
			"1 < 1",
			false,
		},
		{
			"1 > 1",
			false,
		},
		{
			"1 == 1",
			true,
		},
		{
			"1 != 1",
			false,
		},
		{
			"1 == 2",
			false,
		},
		{
			"1 != 2",
			true,
		},
		{
			"true == true",
			true,
		},
		{
			"false == false",
			true,
		},
		{
			"true == false",
			false,
		},
		{
			"true != false",
			true,
		},
		{
			"false != true",
			true,
		},
		{
			"(1 < 2) == true",
			true,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; }"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function but %T", evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("wanted 1 parameter but got %d", len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("should've gotten x as a parameter but got %q", fn.Parameters[0])
	}

	expectedBody := "{ (x + 2) }"

	if fn.Body.String() != expectedBody {
		t.Fatalf("should've gotten %q but got %q for the body", expectedBody, fn.Body.String())
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"if (true) { 10 }",
			10,
		},
		{
			"if (false) { 10 }",
			nil,
		},
		{
			"if (1) { 10 }",
			10,
		},
		{
			"if (1 < 2) { 10 }",
			10,
		},
		{
			"if (1 > 2) { 10 }",
			nil,
		},
		{
			"if (1 > 2) { 10 } else { 20 }",
			20,
		},
		{
			"if (1 < 2) { 10 } else { 20 }",
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			if !testIntegerObject(t, evaluated, int64(integer)) {
				t.Errorf("%s", tt.input)
			}
		} else {
			if !testNullObject(t, evaluated) {
				t.Errorf("%s", tt.input)
			}
		}
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

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"return 10;",
			10,
		},
		{
			"return 10; 9;",
			10,
		},
		{
			"return 2 * 5; 9;",
			10,
		},
		{
			"9; return 2 * 9; 9;",
			18,
		},
		{
			`if (10 > 1) {
				if (10 > 1) {
					return 10;
				}

				return 1
			}`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestStringConcat(t *testing.T) {
	input := `"Hello" + " " + "world!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello world!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello, World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("I wanted a String but I got a %T", evaluated)
	}

	if str.Value != "Hello, World!" {
		t.Errorf("String has the wrong value")
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"let a = 5; a;",
			5,
		},
		{
			"let a = 5 * 5; a;",
			25,
		},
		{
			"let a = 5; let b = a; b;",
			5,
		},
		{
			"let a = 5; let b = a; let c = a + b + 5; c;",
			15,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func testEval(input string) object.Object {

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
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

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL but (%+v)", obj)
		return false
	}

	return true
}
