package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

// Define constants for the Boolean object
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

// Eval is a function that evaluates an AST node
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	}

	return nil
}

// evalStatements is a helper function that takes in a slice of statements and
// evaluates each statement in the slice
func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

// nativeBoolToBooleanObject is a helper function that takes in a boolean and
// returns a pointer to a Boolean object
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// evalPrefixExpression is a helper function that takes in an operator and an
// object and evaluates the prefix expression
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)

	default:
		return NULL
	}
}

// evalInfixExpression is a helper function that takes in an operator and two
// objects and evaluates the infix expression
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	// Integer expressions
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	default:
		return NULL
	}
}

// evalIntegerInfixExpression is a helper function that takes in an operator and
// two objects and evaluates the infix expression
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	// Get the values from the objects
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	// Perform the operation
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}

	case "-":
		return &object.Integer{Value: leftVal - rightVal}

	case "*":
		return &object.Integer{Value: leftVal * rightVal}

	case "/":
		return &object.Integer{Value: leftVal / rightVal}

	default:
		return NULL
	}
}

// evalMinusPrefixOperatorExpression is a helper function that takes in an object
// and evaluates the minus prefix operator expression
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	// Check if the object is an integer
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	// Get the value from the object
	value := right.(*object.Integer).Value

	// Perform the operation
	return &object.Integer{Value: -value}
}

// evalBangOperatorExpression is a helper function that takes in an object and
// evaluates the bang operator expression
func evalBangOperatorExpression(right object.Object) object.Object {
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
