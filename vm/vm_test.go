// vm/vm_test.go

package vm

import (
	"fmt"
	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

// TestTensorMath is a function to test the tensor math bits
func TestTensorMath(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let x = tensor([1], [5.0]) - tensor([1], [4.0]);
			x;`,
			expected: object.Tensor{Shape: []int64{1}, Data: []float64{1.0}},
		},
		{
			input: `
			let x = tensor([1], [5.0]) + tensor([1], [4.0]);
			x;`,
			expected: object.Tensor{Shape: []int64{1}, Data: []float64{9.0}},
		},
		{
			input: `
			let x = tensor([1], [5.0]) * tensor([1], [4.0]);
			x;`,
			expected: object.Tensor{Shape: []int64{1}, Data: []float64{20.0}},
		},
		{
			input: `
			let x = tensor([1], [5.0]) / tensor([1], [2.5]);
			x;`,
			expected: object.Tensor{Shape: []int64{1}, Data: []float64{2.0}},
		},
	}

	runVmTests(t, tests)
}

// TestTensorLiteral is a function to test the tensor literal bits
func TestTensorLiteral(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let x = tensor([3],[1.0,2.0,3.0]);
			x;
			`,
			expected: Null,
		},
	}

	runVmTests(t, tests)
}

// TestFloatLiteral is a function to test the float literal bits
func TestFloatLiteral(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let x = 5.2 + 10.1;
			x;
			`,
			expected: 15.3,
		},
	}

	runVmTests(t, tests)
}

// TestImportLiteral is a function to test the import literal
func TestImportLiteral(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			import "../helper.mky";
			sigmoid(0.5);
			`,
			expected: 0.5,
		},
		{
			input: `
			import "../helper.mky";
			sigmoid(0.4);
			`,
			expected: 0.4,
		},
	}

	runVmTests(t, tests)
}

// TestRecursiveFunctions
func TestRecursiveFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let countDown = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					countDown(x - 1);
				}
			};
			countDown(1);
			`,
			expected: 0,
		},
		{
			input: `
			let countDown = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					countDown(x - 1);
				}
			};
			countDown(2);
			`,
			expected: 0,
		},
		{
			input: `
			let wrapper = fn() {
				let countDown = fn(x) {
					if (x == 0) {
						return 0;
					} else {
						countDown(x - 1);
					}
				};
				countDown(2);
			};
			wrapper();
			`,
			expected: 0,
		},
	}

	runVmTests(t, tests)

}

// TestClosures
func TestClosures(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let newClosure = fn(a) {
				fn() { a; };
			};
			let closure = newClosure(99);
			closure();
			`,
			expected: 99,
		},
		{
			input: `
				let newAdder = fn(a, b) {
					fn(c) { a + b + c };
				};
				let adder = newAdder(1, 2);
				adder(8);
				`,
			expected: 11,
		},
		{
			input: `
				let newAdder = fn(a, b) {
					let c = a + b;
					fn(d) { c + d };
				};
				let adder = newAdder(1, 2);
				adder(8);
				`,
			expected: 11,
		},
		{
			input: `
				let newAdderOuter = fn(a, b) {
					let c = a + b;
					fn(d) {
						let e = d + c;
						fn(f) { e + f; };
					};
				};
				let newAdderInner = newAdderOuter(1, 2);
				let adder = newAdderInner(3);
				adder(8);
				`,
			expected: 14,
		},
		{
			input: `
				let a = 1;
				let newAdderOuter = fn(b) {
					fn(c) {
						fn(d) { a + b + c + d };
					};
				};
				let newAdderInner = newAdderOuter(2);
				let adder = newAdderInner(3);
				adder(8);
				`,
			expected: 14,
		},
		{
			input: `
				let newClosure = fn(a, b) {
					let one = fn() { a; };
					let two = fn() { b; };
					fn() { one() + two(); };
				};
				let closure = newClosure(9, 90);
				closure();
			`,
			expected: 99,
		},
	}

	runVmTests(t, tests)

}

// TestBuiltinFunctions
func TestBuiltinFunctions(t *testing.T) {
	tests := []vmTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, &object.Error{Message: "argument to `len` not supported, got INTEGER"}},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`puts("hello", "world!")`, Null},
		{`first([1,2,3])`, 1},
		{`first([])`, Null},
		{`first(1)`, &object.Error{Message: "argument to `first` must be ARRAY, got INTEGER"}},
		{`rest([1,2,3])`, []int{2, 3}},
		{`rest([])`, Null},
		{`push([],1)`, []int{1}},
		{`push(1,1)`, &object.Error{Message: "argument to `push` must be ARRAY, got INTEGER"}},
	}

	runVmTests(t, tests)
}

// TestCallingFunctionWithErrors
func TestCallingFunctionWithErrors(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `fn() { 1; }(1);
			`,
			expected: "wrong number of arguments: want=0, got=1",
		},
		{
			input: `fn(a) { a; }();
			`,
			expected: "wrong number of arguments: want=1, got=0",
		},
		{
			input:    `fn(a, b) { a + b; }(1);`,
			expected: "wrong number of arguments: want=2, got=1",
		},
	}

	for i, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		machine := New(comp.Bytecode())

		err = machine.Run()
		if err == nil {
			t.Fatalf("test[%d] - expected error", i)
		}

		if err.Error() != tt.expected {
			t.Fatalf("test[%d] - wrong error message. expected=%q, got=%q",
				i, tt.expected, err.Error())
		}
	}
}

// TestCallingFunctionsWithArgumentsAndBindings is a function to test the calling functions with arguments and bindings
func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `let identity = fn(a) { a; }; identity(4);`,
			expected: 4,
		},
		{
			input:    `let sum = fn(a, b) { a + b; }; sum(1, 2);`,
			expected: 3,
		},
		{
			input: `
			let globalNum = 10;
			
			let sum  = fn(a, b) {
				let c =  a + b;
				c + globalNum;
			};
			
			let outer = fn() {
				sum(1, 2) + sum(3, 4) + globalNum;
			};
			
			outer() + globalNum;
			`,
			expected: 50,
		},
	}

	runVmTests(t, tests)
}

// TestCallingFunctionsWithBindings is a function test function bindings
func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let one = fn() { let one = 1; one };
			one();
			`,
			expected: 1,
		},
		{
			input: `
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
			oneAndTwo();
			`,
			expected: 3,
		},
		{
			input: `
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
			let threeAndFour = fn() { let three = 3; let four = 4; three + four; };
			oneAndTwo() + threeAndFour();
			`,
			expected: 10,
		},
		{
			input: `
			let firstFoobar = fn() { let foobar = 50; foobar; };
			let secondFoobar = fn() { let foobar = 100; foobar; };
			firstFoobar() + secondFoobar();
			`,
			expected: 150,
		},
		{
			input: `
			let globalSeed = 50;
			let minusOne = fn() {
				let num = 1;
				globalSeed - num;
			};
			let minusTwo = fn() {
				let num = 2;
				globalSeed - num;
			};
			minusOne() + minusTwo();
			`,
			expected: 97,
		},
	}

	runVmTests(t, tests)
}

// TestFirstClassFunctions is a function to test the first class functions
func TestFirstClassFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `let returnsOneReturner = fn() { let returnsOne = fn() { 1; }; returnsOne; }; returnsOneReturner()();`,
			expected: 1,
		},
	}

	runVmTests(t, tests)
}

// TestFunctionsWithoutReturnValue is a function to test the functions without return value
func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `let noReturn = fn() { }; noReturn();`,
			expected: Null,
		},
		{
			input:    `let noReturn = fn() { }; let noReturnTwo = fn() { noReturn(); }; noReturn(); noReturnTwo();`,
			expected: Null,
		},
	}

	runVmTests(t, tests)
}

// TestFunctionWithReturnValue is a function to test the function with return value
func TestFunctionWithReturnValue(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `let earlyexit = fn() { return 99; 100; }; earlyexit();`,
			expected: 99,
		},
		{
			input:    `let earlyexit = fn() { return 99; return 100; }; earlyexit();`,
			expected: 99,
		},
	}

	runVmTests(t, tests)
}

// TestFunctionWithoutParameters is a function to test the function without parameters
func TestFunctionWithoutParameters(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `let noArg = fn() { 24; }; noArg();`,
			expected: 24,
		},
		{
			input:    `let noArg = fn() { 5 + 10; }; let noArgReturn = noArg(); noArgReturn;`,
			expected: 15,
		},
		{
			input:    `let noArg = fn() { 5 + 10; }; noArg();`,
			expected: 15,
		},
		{
			input: `let one = fn() { 1; };
					let two = fn() { 2; };
					one() + two();`,
			expected: 3,
		},
		{
			input: `let a = fn() { 1; };
					let b = fn() { a() + 1; };
					let c = fn() { b() + 1; };
					c();`,
			expected: 3,
		},
	}

	runVmTests(t, tests)
}

// TestIndexExpressions is a function to test the index expressions
func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", Null},
		{"[1, 2, 3][99]", Null},
		{"[1][-1]", Null},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", Null},
		{"{}[0]", Null},
	}

	runVmTests(t, tests)
}

// TestHashLiterals is a function to test the hash literals
func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			"{}", map[object.HashKey]int64{},
		},
		{
			"{1: 2, 2: 3}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 2}).HashKey(): 3,
			},
		},
		{
			"{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}

	runVmTests(t, tests)
}

// TestArrayLiterals is a function to test the array literals
func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}

	runVmTests(t, tests)
}

// TestStringExpressions is a function to test the string expressions
func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}

	runVmTests(t, tests)
}

// TestGlobalLetStatements is a function to test the global let statements
func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}

	runVmTests(t, tests)
}

// TestBooleanExpressions is a function to test the boolean expressions
func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) { 5; })", true},
	}

	runVmTests(t, tests)
}

// Test Conditionals is a function to test the conditionals
func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}

	runVmTests(t, tests)
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.LastPoppedStackElem()
		testExpectedObject(t, tt.expected, stackElem)

		for i, constant := range comp.Bytecode().Constants {
			fmt.Printf("CONSTANT %d %p (%T):\n", i, constant, constant)

			switch constant := constant.(type) {
			case *object.CompiledFunction:
				fmt.Printf(" Instructions:\n%s", constant.Instructions)
			case *object.Integer:
				fmt.Printf(" Value: %d\n", constant.Value)
			case *object.Tensor:
				fmt.Printf(" Value: %+v\n", constant)
			}
			fmt.Printf("\n")
		}
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)
		if !ok {
			t.Errorf("object is not Hash. got=%T", actual)
			return
		}

		if len(hash.Pairs) != len(expected) {
			t.Errorf("hash has wrong num of pairs. got=%d", len(hash.Pairs))
			return
		}

		for k, v := range expected {
			pair, ok := hash.Pairs[k]
			if !ok {
				t.Errorf("no pair for given key in Pairs")
			}

			err := testIntegerObject(v, pair.Value)
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object is not Array. got=%T", actual)
			return
		}

		if len(array.Elements) != len(expected) {
			t.Errorf("wrong number of elements. want=%d, got=%d", len(expected), len(array.Elements))
			return
		}

		for i, expectedElem := range expected {
			err := testIntegerObject(int64(expectedElem), array.Elements[i])
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case bool:
		err := testBooleanObject(expected, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case object.Tensor:
		err := testTensorObject(expected, actual)
		if err != nil {
			t.Errorf("testTensorObject failed: %+v", err)
		}
	case *object.Error:
		errObj, ok := actual.(*object.Error)
		if !ok {
			t.Errorf("object is not Error: %T (%+v)", actual, actual)
			return
		}
		if errObj.Message != expected.Message {
			t.Errorf("wrong error message. expected=%q, got=%q", expected.Message, errObj.Message)
		}
	}
}

func testTensorObject(expected object.Tensor, actual object.Object) error {
	_, ok := actual.(*object.Tensor)
	if !ok {
		return fmt.Errorf("object is not a Tensor. got=%T", actual)
	}

	if !shapesEqual(expected.Shape, actual.(*object.Tensor).Shape) {
		return fmt.Errorf("shapes have wrong value got=%+v, want=%+v", actual.(*object.Tensor).Shape, expected.Shape)
	}

	if !dataEqual(expected.Data, actual.(*object.Tensor).Data) {
		return fmt.Errorf("data has wrong value. got=%+v, want=%+v", actual.(*object.Tensor).Data, expected.Data)
	}
	return nil
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T", actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%s, want=%s", result.Value, expected)
	}

	return nil
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

// testBooleanObject is a helper function to test the boolean object
func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T", actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}

	return nil
}

// testIntegerObject is a helper function to test the integer object
func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T", actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}
