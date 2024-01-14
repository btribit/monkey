package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression               // prefixParseFn is a function that returns an expression
	infixParseFn  func(ast.Expression) ast.Expression // infixParseFn is a function that returns an expression
)

const (
	_ int = iota
	// LOWEST is the lowest precedence
	LOWEST
	// EQUALS is the precedence of the equals sign
	EQUALS
	// LESSGREATER is the precedence of the less than and greater than signs
	LESSGREATER
	// SUM is the precedence of the sum sign
	SUM
	// PRODUCT is the precedence of the product sign
	PRODUCT
	// PREFIX is the precedence of the prefix sign
	PREFIX
	// CALL is the precedence of the call sign
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// Parser is a struct that holds the lexer and the currentToken
type Parser struct {
	l      *lexer.Lexer
	errors []string

	currentToken token.Token
	peekToken    token.Token

	prefixParseFns map[token.TokenType]prefixParseFn // prefixParseFns is a map of prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn  // infixParseFns is a map of infixParseFn
}

// New is a function that creates a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn) // Initialize the prefixParseFns
	p.registerPrefix(token.IDENT, p.parseIdentifier)           // Register the parseIdentifier function
	p.registerPrefix(token.INT, p.parseIntegerLiteral)         // Register the parseIntegerLiteral function
	p.registerPrefix(token.BANG, p.parsePrefixExpression)      // Register the parsePrefixExpression function
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)     // Register the parsePrefixExpression function

	p.infixParseFns = make(map[token.TokenType]infixParseFn) // Initialize the infixParseFns
	p.registerInfix(token.PLUS, p.parseInfixExpression)      // Register the parseInfixExpression function
	p.registerInfix(token.MINUS, p.parseInfixExpression)     // Register the parseInfixExpression function
	p.registerInfix(token.SLASH, p.parseInfixExpression)     // Register the parseInfixExpression function
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)  // Register the parseInfixExpression function
	p.registerInfix(token.EQ, p.parseInfixExpression)        // Register the parseInfixExpression function
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)    // Register the parseInfixExpression function
	p.registerInfix(token.LT, p.parseInfixExpression)        // Register the parseInfixExpression function
	p.registerInfix(token.GT, p.parseInfixExpression)        // Register the parseInfixExpression function

	// Read two tokens so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken, // Set the token
		Operator: p.currentToken.Literal,
	}

	p.nextToken() // Advance the current token

	expression.Right = p.parseExpression(PREFIX) // Parse the right expression

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken, // Set the token
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence() // Get the precedence of the current token

	p.nextToken() // Advance the current token

	expression.Right = p.parseExpression(precedence) // Parse the right expression

	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal} // Create a new identifier
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn // Register the prefixParseFn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn // Register the infixParseFn
}

// nextToken is a helper function that advances both currentToken and peekToken
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} // Create a new program

	program.Statements = []ast.Statement{} // Initialize the statements

	// Loop through all the statements
	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement() // Parse the statement
		if stmt != nil {
			program.Statements = append(program.Statements, stmt) // Append the statement to the program
		}
		p.nextToken()
	}

	return program
}

// parseStatement is a helper function that parses a statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement() // parseLetStatement is a helper function
	case token.RETURN:
		return p.parseReturnStatement() // parseReturnStatement is a helper function
	default:
		return p.parseExpressionStatement() // parseExpressionStatement is a helper function
	}
}

// parseLetStatement is a helper function that parses a let statement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken} // Create a new let statement

	// Check if the next token is an identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal} // Set the identifier

	// Check if the next token is an equal sign
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We're skipping the expressions until we encounter a semicolon
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement is a helper function that parses a return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken} // Create a new return statement

	p.nextToken() // Advance the current token

	// TODO: We're skipping the expressions until we encounter a semicolon
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement is a helper function that parses an expression statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken} // Create a new expression statement

	stmt.Expression = p.parseExpression(LOWEST) // Parse the expression

	// Check if the next token is a semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken() // Advance the current token
	}

	return stmt
}

// parseExpression is a helper function that parses an expression
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type] // Get the prefixParseFn
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type) // Add an error to the errors slice
		return nil
	}
	leftExp := prefix() // Parse the prefix expression

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type] // Get the infixParseFn
		if infix == nil {
			return leftExp
		}

		p.nextToken() // Advance the current token

		leftExp = infix(leftExp) // Parse the infix expression
	}

	// Loop through the infixParseFns and check if the precedence is higher than the precedence of the current token
	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg) // Add an error to the errors slice
}

// parseIntegerLiteral is a helper function that parses an integer literal
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken} // Create a new integer literal

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64) // Convert the literal to an integer
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg) // Add an error to the errors slice
		return nil
	}

	lit.Value = value // Set the value

	return lit
}

// currentTokenIs is a helper function that checks if the current token is of a certain type
func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// peekTokenIs is a helper function that checks if the peek token is of a certain type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek is a helper function that checks if the peek token is of a certain type
// If it is, it advances both currentToken and peekToken
func (p *Parser) expectPeek(t token.TokenType) bool {
	// Check if the peek token is of the correct type
	if p.peekTokenIs(t) {
		p.nextToken() // Advance both currentToken and peekToken
		return true
	} else {
		p.peekError(t) // Add an error to the errors slice
		return false
	}
}

// Errors is a function that returns the errors
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError is a helper function that adds an error to the errors slice
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// peekPrecedence is a helper function that returns the precedence of the peek token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok { // Check if the peek token is in the precedences map
		return p
	}

	return LOWEST
}

// currentPrecedence is a helper function that returns the precedence of the current token
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok { // Check if the current token is in the precedences map
		return p
	}

	return LOWEST
}
