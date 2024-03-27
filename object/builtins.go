// object/builtins

package object

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

func random() float64 {
	rand.Seed(time.Now().UnixNano()) // Initialize the global random number generator
	return rand.Float64()
}

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
		},
	},
	{
		"puts",
		&Builtin{Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return nil
		},
		},
	},
	{
		"first",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return nil
		},
		},
	},
	{
		"last",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return nil
		},
		},
	},
	{
		"rest",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &Array{Elements: newElements}
			}

			return nil
		},
		},
	},
	{
		"push",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*Array)
			arr.Elements = append(arr.Elements, args[1])

			return arr
		},
		},
	},
	{
		"pop",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `pop` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length > 0 {
				lastElement := arr.Elements[length-1]
				arr.Elements = arr.Elements[:length-1]
				return lastElement
			}

			return nil
		},
		},
	},
	{
		"join",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("first argument to `join` must be ARRAY, got %s", args[0].Type())
			}
			if args[1].Type() != STRING_OBJ {
				return newError("second argument to `join` must be STRING, got %s", args[1].Type())
			}

			arr := args[0].(*Array)
			sep := args[1].(*String)

			strs := make([]string, len(arr.Elements))
			for i, obj := range arr.Elements {
				strs[i] = obj.Inspect()
			}

			return &String{Value: strings.Join(strs, sep.Value)}
		},
		},
	},
	{
		"random",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) > 0 {
				return newError("random() takes no arguments")
			}
			return &Float{Value: random()}
		},
		},
	},
	{
		"exp",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. exp() requires exactly one argument.")
			}

			// Check for the argument type
			switch arg := args[0].(type) {
			case *Float:
				return &Float{Value: math.Exp(arg.Value)}
			case *Integer:
				return &Float{Value: math.Exp(float64(arg.Value))}
			default:
				return newError("argument to `exp` must be a number")
			}

		},
		},
	},
}

// newError returns a new error object with the given format and arguments.
func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

// GetBuiltInByName returns the built-in object with the given name.
func GetBuiltInByName(name string) *Builtin {
	for _, bi := range Builtins {
		if bi.Name == name {
			return bi.Builtin
		}
	}
	return nil
}
