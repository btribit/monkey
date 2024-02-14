package evaluator

import (
	"monkey/object"
	"testing"
)

// testIncorrectParamCount tests for built-in functions with incorrect parameter counts
func TestIncorrectParamCount(t *testing.T) {
	tests := []struct {
		input          string
		expectedError  bool
		expectedErrMsg string
	}{
		{`len("hello", "world")`, true, "wrong number of arguments. got=2, want=1"},
		{`first([1, 2, 3], "extra")`, true, "wrong number of arguments. got=2, want=1"},
		{`last([], "extra", "extra")`, true, "wrong number of arguments. got=3, want=1"},
		{`rest()`, true, "wrong number of arguments. got=0, want=1"},
		{`push([1, 2, 3], 4, 5)`, true, "wrong number of arguments. got=3, want=2"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if !tt.expectedError {
			t.Errorf("expected no error, got %T", evaluated)
			continue
		}

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("expected error object, got %T", evaluated)
			continue
		}

		if errObj.Message != tt.expectedErrMsg {
			t.Errorf("expected error message %q, got %q", tt.expectedErrMsg, errObj.Message)
		}
	}
}

func TestBuiltins(t *testing.T) {
	tests := []struct {
		input          string
		expectedType   object.ObjectType
		expectedValue  interface{}
		expectedError  bool
		expectedErrMsg string
	}{
		// len tests
		{`len("hello")`, object.INTEGER_OBJ, 5, false, ""},
		{`len("")`, object.INTEGER_OBJ, 0, false, ""},
		{`len([1, 2, 3])`, object.INTEGER_OBJ, 3, false, ""},
		{`len([])`, object.INTEGER_OBJ, 0, false, ""},
		{`len(123)`, object.ERROR_OBJ, nil, true, "argument to `len` not supported, got INTEGER"},

		// first tests
		{`first([1, 2, 3])`, object.INTEGER_OBJ, 1, false, ""},
		{`first([])`, object.NULL_OBJ, nil, false, ""},
		{`first("hello")`, object.ERROR_OBJ, nil, true, "argument to `first` must be ARRAY, got STRING"},

		// last tests
		{`last([1, 2, 3])`, object.INTEGER_OBJ, 3, false, ""},
		{`last([])`, object.NULL_OBJ, nil, false, ""},
		{`last("hello")`, object.ERROR_OBJ, nil, true, "argument to `last` must be ARRAY, got STRING"},

		// rest tests
		{`rest([1, 2, 3])`, object.ARRAY_OBJ, []int{2, 3}, false, ""},
		{`rest([])`, object.NULL_OBJ, nil, false, ""},
		{`rest("hello")`, object.ERROR_OBJ, nil, true, "argument to `rest` must be ARRAY, got STRING"},

		// push tests
		{`push([1, 2], 3)`, object.ARRAY_OBJ, []int{1, 2, 3}, false, ""},
		{`push([], 1)`, object.ARRAY_OBJ, []int{1}, false, ""},
		{`push("hello", 1)`, object.ERROR_OBJ, nil, true, "argument to `push` must be ARRAY, got STRING"},

		// puts tests
		{`puts("hello", 123, [1, 2, 3])`, object.NULL_OBJ, nil, false, ""},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if tt.expectedError {
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("expected error object, got %T", evaluated)
				continue
			}

			if errObj.Message != tt.expectedErrMsg {
				t.Errorf("expected error message %q, got %q", tt.expectedErrMsg, errObj.Message)
			}
		} else {
			testObject(t, evaluated, tt.expectedType, tt.expectedValue)
		}
	}
}

// testObject is a helper function to test the object returned from the evaluator
func testObject(t *testing.T, obj object.Object, expectedType object.ObjectType, expectedValue interface{}) {
	if obj.Type() != expectedType {
		t.Errorf("object has wrong type. expected=%q, got=%q", expectedType, obj.Type())
	}

	switch expectedType {
	case object.INTEGER_OBJ:
		val, ok := expectedValue.(int)
		if !ok {
			t.Fatalf("expectedValue is not an int. got=%T", expectedValue)
		}

		testIntegerObject(t, obj, int64(val))
	case object.BOOLEAN_OBJ:
		val, ok := expectedValue.(bool)
		if !ok {
			t.Fatalf("expectedValue is not a bool. got=%T", expectedValue)
		}

		testBooleanObject(t, obj, val)
	case object.NULL_OBJ:
		testNullObject(t, obj)
	case object.ARRAY_OBJ:
		val, ok := expectedValue.([]int)
		if !ok {
			t.Fatalf("expectedValue is not a []int. got=%T", expectedValue)
		}

		testArrayObject(t, obj, val)
	case object.ERROR_OBJ:
		val, ok := expectedValue.(string)
		if !ok {
			t.Fatalf("expectedValue is not a string. got=%T", expectedValue)
		}

		testErrorObject(t, obj, val)
	}
}

// testArrayObject is a helper function to test the array object returned from the evaluator
func testArrayObject(t *testing.T, obj object.Object, expected []int) {
	result, ok := obj.(*object.Array)
	if !ok {
		t.Fatalf("object is not an Array. got=%T", obj)
	}

	if len(result.Elements) != len(expected) {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	for i, el := range expected {
		testIntegerObject(t, result.Elements[i], int64(el))
	}
}

// testErrorObject is a helper function to test the error object returned from the evaluator
func testErrorObject(t *testing.T, obj object.Object, expected string) {
	errObj, ok := obj.(*object.Error)
	if !ok {
		t.Fatalf("object is not an Error. got=%T", obj)
	}

	if errObj.Message != expected {
		t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
	}
}
