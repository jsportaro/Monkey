package parser

import (
	"Monkey/ast"
	"Monkey/lexer"
	"fmt"
	"testing"
)

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

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
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
		if !testIntegerLiteral(t, expression.Right, tt.integerValue) {
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
