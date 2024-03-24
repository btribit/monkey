package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"os"
)

// Define constants for the Boolean object
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

// Eval is a function that evaluates an AST node
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.ImportLiteral:
		return evalImportLiteral(node, env)
	}

	return nil
}

// evalImportLiteral is a helper function that takes in an import literal and an
// environment and evaluates the import literal
func evalImportLiteral(node *ast.ImportLiteral, env *object.Environment) object.Object {
	// Read the file from the node.path into a string
	fileContent, err := os.ReadFile(node.Path)
	if err != nil {
		return newError("error reading import file: %s", err.Error())
	}

	// convert fileContent to a string
	fileContentString := string(fileContent)

	l := lexer.New(fileContentString)
	p := parser.New(l)
	program := p.ParseProgram()
	evaluated := Eval(program, env)
	if evaluated != nil {
		return evaluated
	}

	return NULL
}

// evalHashLiteral is a helper function that takes in a hash literal and an
// environment and evaluates the hash literal
func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

// evalIndexExpression is a helper function that takes in two objects and
// evaluates the index expression
func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

// evalHashIndexExpression is a helper function that takes in two objects and
// evaluates the hash index expression
func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

// evalArrayIndexExpression is a helper function that takes in two objects and
// evaluates the array index expression
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

// applyFunction is a helper function that takes in a function and a slice of
// arguments and applies the function to the arguments
func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch function := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(function, args)
		evaluated := Eval(function.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		if result := function.Fn(args...); result != nil {
			return result
		}
		return NULL

	default:
		return newError("not a function: %s", fn.Type())
	}
}

// extendFunctionEnv is a helper function that takes in a function and a slice of
// arguments and extends the function's environment with the arguments
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// newError is a helper function that takes in a string and returns an error
// object
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// isError is a helper function that takes in an object and returns a boolean
// indicating whether the object is an error
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// evalExpressions is a helper function that takes in a slice of expressions and
// evaluates each expression in the slice
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// evalBlockStatements is a helper function that takes in a slice of statements
// and evaluates each statement in the slice
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		// Check if the result is a return value or an error
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
			return result
		}
	}

	return result
}

// evalProgram is a helper function that takes in a program and evaluates the
// program
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		// Check if the result is a return value
		case *object.ReturnValue:
			return result.Value

		// Check if the result is an error
		case *object.Error:
			return result
		}
	}

	return result
}

// evalIfExpression is a helper function that takes in an if expression and
// evaluates the expression
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	// Evaluate the condition
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	// Check if the condition is true
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false

	case TRUE:
		return true

	case FALSE:
		return false

	default:
		return true
	}
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
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// evalInfixExpression is a helper function that takes in an operator and two
// objects and evaluates the infix expression
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	// Integer expressions
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalStringInfixExpression is a helper function that takes in an operator and
// two objects and evaluates the infix expression
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	// Get the values from the objects
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	// Perform the operation
	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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

	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)

	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)

	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)

	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalFloatInfixExpression is a helper function that takes in an operator and
// two objects and evaluates the infix expression
// to copypasta or not to copypasta
func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	// Get the values from the objects
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	// Perform the operation
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}

	case "-":
		return &object.Float{Value: leftVal - rightVal}

	case "*":
		return &object.Float{Value: leftVal * rightVal}

	case "/":
		return &object.Float{Value: leftVal / rightVal}

	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)

	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)

	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)

	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalMinusPrefixOperatorExpression is a helper function that takes in an object
// and evaluates the minus prefix operator expression
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	// Check if the object is an integer
	if right.Type() != object.INTEGER_OBJ && right.Type() != object.FLOAT_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	if right.Type() == object.FLOAT_OBJ {
		value := right.(*object.Float).Value

		return &object.Float{Value: -value}
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

// evalIdentifier is a helper function that takes in an identifier and evaluates
// the identifier
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}
