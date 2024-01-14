package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

// checkParserErrors is a helper function that checks if there are any parser errors
func checkParserErrors(t *testing.T, p *Parser) {
	// Check if there are any parser errors
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	// Print the errors
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. Got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	// Loop through the tests and check if the identifier is correct
	for i, tt := range tests {
		stmt := program.Statements[i] // Get the statement
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

// testLetStatement is a helper function that checks if the statement is a let statement
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	// Check if the statement is a let statement
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. Got %q", s.TokenLiteral())
		return false
	}

	// Type assertion to get the *ast.LetStatement
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. Got %T", s)
		return false
	}

	// Check if the identifier is correct
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. Got %s", name, letStmt.Name.Value)
		return false
	}

	// Check if the identifier token literal is correct
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. Got %s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

// TestReturnStatements is a function that tests return statements
func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	// Check if the program contains 3 statements
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. Got %d", len(program.Statements))
	}

	// Loop through the statements and check if they are return statements
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement) // Type assertion to get the *ast.ReturnStatement
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. Got %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" { // Check if the token literal is correct
			t.Errorf("returnStmt.TokenLiteral not 'return', Got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. Got %d", len(program.Statements))
	}
	// Type assertion to get the *ast.ExpressionStatement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T", program.Statements[0])
	}
	// Type assertion to get the *ast.Identifier
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. Got %T", stmt.Expression)
	}
	if ident.Value != "foobar" { // Check if the value is correct
		t.Errorf("ident.Value not %s. Got %s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" { // Check if the token literal is correct
		t.Errorf("ident.TokenLiteral not %s. Got %s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. Got %d", len(program.Statements))
	}
	// Type assertion to get the *ast.ExpressionStatement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T", program.Statements[0])
	}
	// Type assertion to get the *ast.IntegerLiteral
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. Got %T", stmt.Expression)
	}
	if literal.Value != 5 { // Check if the value is correct
		t.Errorf("literal.Value not %d. Got %d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" { // Check if the token literal is correct
		t.Errorf("literal.TokenLiteral not %s. Got %s", "5", literal.TokenLiteral())
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	// Type assertion to get the *ast.IntegerLiteral
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. Got %T", il)
		return false
	}
	if integ.Value != value { // Check if the value is correct
		t.Errorf("integ.Value not %d. Got %d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) { // Check if the token literal is correct
		t.Errorf("integ.TokenLiteral not %d. Got %s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func TestParsingPrefixExpression(t *testing.T) {
	// Create a struct to represent the test case
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},   // !5
		{"-15;", "-", 15}, // -15
	}

	// Loop through the tests and check if the prefix expression is correct
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p) // Check if there are any parser errors

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. Got %d", len(program.Statements))
		}
		// Type assertion to get the *ast.ExpressionStatement
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T", program.Statements[0])
		}
		// Type assertion to get the *ast.PrefixExpression
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not *ast.PrefixExpression. Got %T", stmt.Expression)
		}
		if exp.Operator != tt.operator { // Check if the operator is correct
			t.Fatalf("exp.Operator is not '%s'. Got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) { // Check if the value is correct
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	// Create a struct to represent the test case
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},   // 5 + 5
		{"5 - 5;", 5, "-", 5},   // 5 - 5
		{"5 * 5;", 5, "*", 5},   // 5 * 5
		{"5 / 5;", 5, "/", 5},   // 5 / 5
		{"5 > 5;", 5, ">", 5},   // 5 > 5
		{"5 < 5;", 5, "<", 5},   // 5 < 5
		{"5 == 5;", 5, "==", 5}, // 5 == 5
		{"5 != 5;", 5, "!=", 5}, // 5 != 5
	}

	// Loop through the tests and check if the infix expression is correct
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p) // Check if there are any parser errors

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. Got %d", len(program.Statements))
		}
		// Type assertion to get the *ast.ExpressionStatement
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T", program.Statements[0])
		}
		// Type assertion to get the *ast.InfixExpression
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt is not *ast.InfixExpression. Got %T", stmt.Expression)
		}
		if !testIntegerLiteral(t, exp.Left, tt.leftValue) { // Check if the left value is correct
			return
		}
		if exp.Operator != tt.operator { // Check if the operator is correct
			t.Fatalf("exp.Operator is not '%s'. Got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.rightValue) { // Check if the right value is correct
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	// Create a struct to represent the test case
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},                                                 // -a * b
		{"!-a", "(!(-a))"},                                                       // !-a
		{"a + b + c", "((a + b) + c)"},                                           // a + b + c
		{"a + b - c", "((a + b) - c)"},                                           // a + b - c
		{"a * b * c", "((a * b) * c)"},                                           // a * b * c
		{"a * b / c", "((a * b) / c)"},                                           // a * b / c
		{"a + b / c", "(a + (b / c))"},                                           // a + b / c
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},             // a + b * c + d / e - f
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},                                   // 3 + 4; -5 * 5
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},                               // 5 > 4 == 3 < 4
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},                               // 5 < 4 != 3 > 4
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"}, // 3 + 4 * 5 == 3 * 1 + 4 * 5
	}

	// Loop through the tests and check if the infix expression is correct
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p) // Check if there are any parser errors

		actual := program.String()
		if actual != tt.expected { // Check if the actual is correct
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
