package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"fmt"
	"testing"
)

func TestBooleanExpression(t *testing.T) {
	infixTests := []struct {
		input string
		value bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserError(t, p)

		stmtCount := len(program.Statements)
		if stmtCount != 1 {
			t.Fatalf("program.Statements does not contain %d statements but %d", 1, stmtCount)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement but %T", program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("expression is not ast.BooleanExpression but %T", stmt.Expression)
		}

		if tt.value != boolean.Value {
			t.Fatalf("boolean was not %T but was %T", tt.value, boolean.Value)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not a ast.ExpressionStatement but was %T", program.Statements[0])
	}

	callExpression, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression wasn't a CallExpression but was %T", stmt.Expression)
	}

	if !testIdentifier(t, callExpression.Function, "add") {
		return
	}

	if len(callExpression.Arguments) != 3 {
		t.Fatalf("Was expected 3 arguments but got %d instead", len(callExpression.Arguments))
	}

	testLiteralExpression(t, callExpression.Arguments[0], 1)
	testInfixExpression(t, callExpression.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExpression.Arguments[2], 4, "+", 5)
}

func TestWhileExpressionParsing(t *testing.T) {
	input := `while (x < y) { let x = x + 1; }`
	program := parseProgram(input, t)

	stmtCount := len(program.Statements)
	if stmtCount != 1 {
		t.Fatalf("program.Statements does not have 2 statement but %d\n", stmtCount)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement but %T", program.Statements[0])
	}

	expression, ok := stmt.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("stmt.Statement is not ast.WhileExpression but %T", stmt.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Body.Statements) != 1 {
		t.Errorf("body is not 1 statement")
	}

	_, ok = expression.Body.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("Statement[0] is not ast.ExpressionStatement")
	}

}

func TestParsingEmptyHashLiteralString(t *testing.T) {
	input := "{}"
	program := parseProgram(input, t)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("wanted a hash but ended up with a %T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs was wrong length.  got %d", len(hash.Pairs))
	}
}

func TestParsingHashLiterals(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	program := parseProgram(input, t)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("wanted a hash but ended up with a %T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs was wrong length.  got %d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral but %T", key)
		}

		expectedValue := expected[literal.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not a ast.ExpressionStatement but was %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression not *ast.Identifier but was %T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Fatalf("expected value 'foobar' but got %s", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("expected TokenLiteral 'foobar' but got %d", ident.TokenLiteral())
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	stmtCount := len(program.Statements)

	if stmtCount != 1 {
		t.Fatalf("program doesn't have 1 statement but %d", stmtCount)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] wasn't a ast.ExpressionStatement but %T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression wasn't a ast.FunctionLiteral but %T", stmt.Expression)
	}

	paramCount := len(function.Parameters)
	if paramCount != 2 {
		t.Fatalf("Expected 2 parameters but found %d", paramCount)
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	functionStmtCount := len(function.Body.Statements)

	if functionStmtCount != 1 {
		t.Fatalf("Expected 1 statment but got %d", functionStmtCount)
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("funciton.Body.Statements[0] wasn't a ast.ExpressionStatement but %T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{
			input:          "fn() {};",
			expectedParams: []string{},
		},
		{
			input:          "fn(x) {};",
			expectedParams: []string{"x"},
		},
		{
			input:          "fn(x,y ,z) {};",
			expectedParams: []string{"x", "y", "z"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserError(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		actualParamCount := len(function.Parameters)
		expectedParamCount := len(tt.expectedParams)

		if actualParamCount != expectedParamCount {
			t.Errorf("Was expecting %d params but got %d", expectedParamCount, actualParamCount)
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	stmtCount := len(program.Statements)
	if stmtCount != 1 {
		t.Fatalf("program.Statements does not have 1 statement but %d\n", stmtCount)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement but %T", program.Statements[0])
	}

	expression, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Statement is not ast.IfExpression but %T", stmt.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement")
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement[0] is not ast.ExpressionStatement")
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("expression.Alternative.Statements was not null")
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	stmtCount := len(program.Statements)
	if stmtCount != 1 {
		t.Fatalf("program.Statements does not have 1 statement but %d\n", stmtCount)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement but %T", program.Statements[0])
	}

	expression, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Statement is not ast.IfExpression but %T", stmt.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement")
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement[0] is not ast.ExpressionStatement")
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative == nil {
		t.Errorf("expression.Alternative.Statements was null")
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Fatalf("alternative is not 1 statement")
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement[0] is not ast.ExpressionStatement")
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	stmtCount := len(program.Statements)

	if stmtCount != 1 {
		t.Fatalf("program has unexpecte number of statements.  was %d", stmtCount)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. was %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression not *ast.Integeriteral. was %T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value was not 5, but %d", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() not '5' but %s", literal.TokenLiteral())
	}

}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{
			"let x = 5;",
			"x",
			5,
		},
		{
			"let y = true",
			"y",
			true,
		},
		{
			"let foobar  = y",
			"foobar",
			"y",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserError(t, p)

		stmtCount := len(program.Statements)

		if stmtCount != 1 {
			t.Fatalf("Wanted 1 statement but got %d", stmtCount)
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value

		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a  * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserError(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected %q but got %q", tt.expected, actual)
		}
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 4, 3 + 3]"
	program := parseProgram(input, t)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)

	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3 but %d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 4)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"
	program := parseProgram(input, t)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExpression, ok := stmt.Expression.(*ast.IndexExpression)

	if !ok {
		t.Fatalf("expression not *ast.IndexExpression but was %T", stmt.Expression)
	}

	if !testIdentifier(t, indexExpression.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExpression.Index, 1, "+", 1) {
		return
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserError(t, p)

		stmtCount := len(program.Statements)
		if stmtCount != 1 {
			t.Fatalf("program.Statements does not contain %d statements but %d", 1, stmtCount)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement but %T", program.Statements[0])
		}

		expression, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expression is not asat.InfixExpression but %T", stmt.Expression)
		}

		if !testLiteralExpression(t, expression.Left, tt.leftValue) {
			t.Errorf("On %s", tt.input)
			return
		}

		if expression.Operator != tt.operator {
			t.Fatalf("expression.Operator is not %s but %s", tt.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, tt.rightValue) {
			t.Errorf("On %s", tt.input)
			return
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserError(t, p)

		stmtCount := len(program.Statements)
		if stmtCount != 1 {
			t.Fatalf("program.Statments does not contain 1 statement but %d\n", stmtCount)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] isn't an ExpressionStatement but %T", program.Statements[0])
		}

		expression, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not an ast.PrefixExpressino, but %T", stmt.Expression)
		}
		if expression.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s but %s", tt.operator, expression.Operator)
		}
		if !testLiteralExpression(t, expression.Right, tt.value) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{
			"return 5;",
			5,
		},
		{
			"return true",
			true,
		},
		{
			"return y",
			"y",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserError(t, p)

		stmtCount := len(program.Statements)

		if stmtCount != 1 {
			t.Fatalf("Wanted 1 statement but got %d", stmtCount)
		}

		stmt := program.Statements[0]

		val := stmt.(*ast.ReturnStatement).ReturnValue

		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"Hello, World!"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)

	if !ok {
		t.Fatalf("wasn't a *ast.StringLiteral but %T", stmt.Expression)
	}

	if literal.Value != "Hello, World!" {
		t.Errorf("I was expecting %q but somehow ended up with %q", input, literal.Value)
	}
}

func checkParserError(t *testing.T, p *Parser) {
	errors := p.Errors()
	errorCount := len(errors)

	if errorCount == 0 {
		return
	}

	t.Errorf("Parser has %d errors", errorCount)

	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}

	t.FailNow()
}

func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) bool {
	bo, ok := expression.(*ast.Boolean)
	if !ok {
		t.Errorf("expression was not *ast.Boolean but was %T", expression)
		return false
	}

	if bo.Value != value {
		t.Errorf("boolean value was not %T but was %T", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	identifier, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expression not *ast.Identifier but %T", expression)
		return false
	}

	if identifier.Value != value {
		t.Errorf("identifier.Value was not %s but was %s", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral() was not %s but was %s", value, identifier.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(
	t *testing.T,
	expression ast.Expression,
	left interface{},
	operator string,
	right interface{}) bool {

	operationExpression, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expression was not ast.OperatorExpression")
		return false
	}

	if !testLiteralExpression(t, operationExpression.Left, left) {
		return false
	}

	if operationExpression.Operator != operator {
		t.Errorf("expression.Operator is not %s but was %q", operator, operationExpression.Operator)
	}

	if !testLiteralExpression(t, operationExpression.Right, right) {
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got %T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d but %d", value, integer.Value)
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() not %d but %s", value, integer.TokenLiteral())
	}

	return true
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got %q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement. got %T", stmt)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got %s", name, letStmt.Name)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expression, int64(v))
	case int64:
		return testIntegerLiteral(t, expression, v)
	case string:
		return testIdentifier(t, expression, v)
	case bool:
		return testBooleanLiteral(t, expression, v)
	}

	t.Errorf("type of exp not handled")

	return false
}

func parseProgram(input string, t *testing.T) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	return program
}
