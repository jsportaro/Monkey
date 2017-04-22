package parser

import (
	"Monkey/ast"
	"Monkey/lexer"
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

func TestLetStatement(t *testing.T) {
	input := `
	
	let x = 5;
	let y = 10;
	let foobar = 8883;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserError(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 999333222;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserError(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("stmt not a ReturnStatement.  Was %T", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral was not 'return'. Was %q", returnStmt.TokenLiteral())
		}
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
