package parser

import (
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/token"
	"fmt"
)

//Parser Make me my AST
type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	currentToken token.Token
	peekedToken  token.Token
}

//New Get me a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []string{}}

	// Clever idea: Prime to parser so I don't need
	// a "current" integer.  Works because lexer in this
	// language implements the enumerator pattern
	p.nextToken()
	p.nextToken()

	return p
}

//ParseProgram Make the AST
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

//Errors Get all errors
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekedToken
	p.peekedToken = p.lexer.NextToken()
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, but was %s instead", t, p.peekedToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	letStmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return letStmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStmt := &ast.ReturnStatement{Token: p.currentToken}

	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return returnStmt
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekedTokenIs(t token.TokenType) bool {
	return p.peekedToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekedTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}
