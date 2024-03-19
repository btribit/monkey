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
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},             // let x = 5;
		{"let y = true;", "y", true},       // let y = true;
		{"let foobar = y;", "foobar", "y"}, // let foobar = y;
	}

	// Loop through the tests and check if the let statements are correct
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p) // Check if there are any parser errors

		// Check if the program contains 1 statement
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
		}

		stmt := program.Statements[0]                          // Get the statement
		if !testLetStatement(t, stmt, tt.expectedIdentifier) { // Check if the statement is a let statement
			return
		}

		val := stmt.(*ast.LetStatement).Value                 // Get the value
		if !testLiteralExpression(t, val, tt.expectedValue) { // Check if the value is correct
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
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},             // return 5;
		{"return true;", true},       // return true;
		{"return foobar;", "foobar"}, // return foobar;
	}

	// Loop through the tests and check if the return statements are correct
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p) // Check if there are any parser errors

		// Check if the program contains 1 statement
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
		}

		stmt := program.Statements[0]      // Get the statement
		if !testReturnStatement(t, stmt) { // Check if the statement is a return statement
			return
		}
		val := stmt.(*ast.ReturnStatement).ReturnValue        // Get the value
		if !testLiteralExpression(t, val, tt.expectedValue) { // Check if the value is correct
			return
		}
	}
}

// testReturnStatement is a helper function that checks if the statement is a return statement
func testReturnStatement(t *testing.T, s ast.Statement) bool {
	// Check if the statement is a return statement
	if s.TokenLiteral() != "return" {
		t.Errorf("s.TokenLiteral not 'return'. Got %q", s.TokenLiteral())
		return false
	}

	// Type assertion to get the *ast.ReturnStatement
	returnStmt, ok := s.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("s not *ast.ReturnStatement. Got %T", s)
		return false
	}

	// Check literal expression
	if returnStmt.TokenLiteral() != "return" {
		t.Errorf("returnStmt.TokenLiteral() not 'return', Got %q", returnStmt.TokenLiteral())
		return false
	}

	return true
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
		integerValue interface{}
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
		if !testLiteralExpression(t, exp.Right, tt.integerValue) { // Check if the value is correct
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	// Create a struct to represent the test case
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},                  // 5 + 5
		{"5 - 5;", 5, "-", 5},                  // 5 - 5
		{"5 * 5;", 5, "*", 5},                  // 5 * 5
		{"5 / 5;", 5, "/", 5},                  // 5 / 5
		{"5 > 5;", 5, ">", 5},                  // 5 > 5
		{"5 < 5;", 5, "<", 5},                  // 5 < 5
		{"5 == 5;", 5, "==", 5},                // 5 == 5
		{"5 != 5;", 5, "!=", 5},                // 5 != 5
		{"true == true", true, "==", true},     // true == true
		{"true != false", true, "!=", false},   // true != false
		{"false == false", false, "==", false}, // false == false
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
		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) { // Check if the infix expression is correct
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
		{"true", "true"},                                                         // true
		{"false", "false"},                                                       // false
		{"3 > 5 == false", "((3 > 5) == false)"},                                 // 3 > 5 == false
		{"3 < 5 == true", "((3 < 5) == true)"},                                   // 3 < 5 == true
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},                               // 1 + (2 + 3) + 4
		{"(5 + 5) * 2", "((5 + 5) * 2)"},                                         // (5 + 5) * 2
		{"2 / (5 + 5)", "(2 / (5 + 5))"},                                         // 2 / (5 + 5)
		{"-(5 + 5)", "(-(5 + 5))"},                                               // -(5 + 5)
		{"!(true == true)", "(!(true == true))"},                                 // !(true == true)
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},                       // a + add(b * c) + d
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"}, // add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},                           // add(a + b + c * d / f + g)
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},                           // a * [1, 2, 3, 4][b * c] * d
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},           // add(a * b[2], b[1], 2 * [1, 2][1])
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
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	// Type assertion to get the *ast.Identifier
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. Got %T", exp)
		return false
	}
	if ident.Value != value { // Check if the value is correct
		t.Errorf("ident.Value not %s. Got %s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value { // Check if the token literal is correct
		t.Errorf("ident.TokenLiteral not %s. Got %s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	// Check if the expression is an integer literal
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	// If the expression is not an integer literal, print an error
	t.Errorf("type of exp not handled. Got=%T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	// Type assertion to get the *ast.Boolean
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. Got %T", exp)
		return false
	}
	if boolean.Value != value { // Check if the value is correct
		t.Errorf("boolean.Value not %t. Got %t", value, boolean.Value)
		return false
	}
	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) { // Check if the token literal is correct
		t.Errorf("boolean.TokenLiteral not %t. Got %s", value, boolean.TokenLiteral())
		return false
	}
	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	// Type assertion to get the *ast.InfixExpression
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. Got %T(%s)", exp, exp)
		return false
	}
	// Check if the left expression is correct
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	// Check if the operator is correct
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. Got %q", operator, opExp.Operator)
		return false
	}
	// Check if the right expression is correct
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func TestBooleanExpression(t *testing.T) {
	// Create a struct to represent the test case
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},   // true
		{"false;", false}, // false
	}

	// Loop through the tests and check if the boolean expression is correct
	for _, tt := range tests {
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
		// Type assertion to get the *ast.Boolean
		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp is not *ast.Boolean. Got %T", stmt.Expression)
		}
		if boolean.Value != tt.expected { // Check if the value is correct
			t.Errorf("boolean.Value not %t. Got %t", tt.expected, boolean.Value)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	// Check if the program contains 1 statement
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
	}

	// Type assertion to get the *ast.ExpressionStatement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T", program.Statements[0])
	}

	// Type assertion to get the *ast.IfExpression
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. Got %T", stmt.Expression)
	}

	// Check if the condition is correct
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	// Check if the consequence is correct
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. Got %d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. Got %T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") { // Check if the identifier is correct
		return
	}

	// Check if the alternative is nil
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. Got %+v", exp.Alternative)
	}

}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	// Check if the program contains 1 statement
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
	}

	// Type assertion to get the *ast.ExpressionStatement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T", program.Statements[0])
	}

	// Type assertion to get the *ast.IfExpression
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. Got %T", stmt.Expression)
	}

	// Check if the condition is correct
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	// Check if the consequence is correct
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. Got %d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. Got %T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") { // Check if the identifier is correct
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. Got %T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") { // Check if the identifier is correct
		return
	}

}

// Test Function Literal Parsing
func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	// Check if the program contains 1 statement
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
	}

	// Type assertion to get the *ast.ExpressionStatement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T", program.Statements[0])
	}

	// Type assertion to get the *ast.FunctionLiteral
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.FunctionLiteral. Got %T", stmt.Expression)
	}

	// Check if the parameters are correct
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. Want 2, Got %d", len(function.Parameters))
	}

	// Check if the first parameter is correct
	testLiteralExpression(t, function.Parameters[0], "x")

	// Check if the second parameter is correct
	testLiteralExpression(t, function.Parameters[1], "y")

	// Check if the body is correct
	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statement. Got %d", len(function.Body.Statements))
	}

	// Type assertion to get the *ast.ExpressionStatement
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. Got %T", function.Body.Statements[0])
	}

	// Check if the body statement is correct
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

// Test Function Parameter Parsing
func TestFunctionParameterParsing(t *testing.T) {
	// Create a struct to represent the test case
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},                     // fn() {}
		{"fn(x) {};", []string{"x"}},                 // fn(x) {}
		{"fn(x, y, z) {};", []string{"x", "y", "z"}}, // fn(x, y, z) {}
	}

	// Loop through the tests and check if the function parameters are correct
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p) // Check if there are any parser errors

		// Check if the program contains 1 statement
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
		}

		// Type assertion to get the *ast.ExpressionStatement
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T", program.Statements[0])
		}

		// Type assertion to get the *ast.FunctionLiteral
		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.FunctionLiteral. Got %T", stmt.Expression)
		}

		// Check if the parameters are correct
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Fatalf("length parameters wrong. Want %d, Got %d", len(tt.expectedParams), len(function.Parameters))
		}

		// Loop through the parameters and check if they are correct
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

// Test Call Expression Parsing
func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	// Check if the program contains 1 statement
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
	}

	// Type assertion to get the *ast.ExpressionStatement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T", program.Statements[0])
	}

	// Type assertion to get the *ast.CallExpression
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.CallExpression. Got %T", stmt.Expression)
	}

	// Check if the function is correct
	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	// Check if the arguments are correct
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. Got %d", len(exp.Arguments))
	}

	// Check if the first argument is correct
	testLiteralExpression(t, exp.Arguments[0], 1)

	// Check if the second argument is correct
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)

	// Check if the third argument is correct
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

// Test String Literal Expression
func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	// Check if the program contains 1 statement
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
	}

	// Type assertion to get the *ast.ExpressionStatement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T", program.Statements[0])
	}

	// Type assertion to get the *ast.StringLiteral
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.StringLiteral. Got %T", stmt.Expression)
	}

	// Check if the value is correct
	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. Got %q", "hello world", literal.Value)
	}
}

// Test Array Literal Parsing
func TestArrayLiteralParsing(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3];"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	// Check if the program contains 1 statement
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
	}

	// Type assertion to get the *ast.ExpressionStatement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T", program.Statements[0])
	}

	// Type assertion to get the *ast.ArrayLiteral
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.ArrayLiteral. Got %T", stmt.Expression)
	}

	// Check if the length of elements is correct
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. Got %d", len(array.Elements))
	}

	// Check if the first element is correct
	testIntegerLiteral(t, array.Elements[0], 1)

	// Check if the second element is correct
	testInfixExpression(t, array.Elements[1], 2, "*", 2)

	// Check if the third element is correct
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

// Test Index Expression Parsing
func TestIndexExpressionParsing(t *testing.T) {
	input := "myArray[1 + 1];"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	// Check if the program contains 1 statement
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
	}

	// Type assertion to get the *ast.ExpressionStatement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T", program.Statements[0])
	}

	// Type assertion to get the *ast.IndexExpression
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IndexExpression. Got %T", stmt.Expression)
	}

	// Check if the left expression is correct
	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	// Check if the index is correct
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

// Test Hash Literal Parsing
func TestHashLiteralParsing(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.HashLiteral. Got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. Got %d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for k, v := range hash.Pairs {
		literal, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not *ast.StringLiteral. Got %T", k)
		}

		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, v, expectedValue)
	}
}

// Test Empty Hash Literal Parsing
func TestEmptyHashLiteralParsing(t *testing.T) {
	input := `{}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.HashLiteral. Got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. Got %d", len(hash.Pairs))
	}
}

// Test Parsing Hash Literals with Expressions
func TestParsingHashLiteralWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.HashLiteral. Got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. Got %d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for k, v := range hash.Pairs {
		literal, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not *ast.StringLiteral. Got %T", k)
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %s found", literal.String())
			continue
		}
		testFunc(v)
	}
}

// TestFunctionLiteralWithName
func TestFunctionLiteralWithName(t *testing.T) {
	input := `let myFunction = fn() { };`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) // Check if there are any parser errors

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statement. Got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.LetStatement. Got=%T", program.Statements[0])
	}

	function, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.FunctionLiteral. Got=%T", stmt.Value)
	}

	if function.Name != "myFunction" {
		t.Errorf("function literal name wrong. Want 'myFunction', Got=%q", function.Name)
	}
}
