// hello.go
package main

import (
	"monkey/object"
)

func Hello(args ...object.Object) object.Object {
	var value string
	value = "Hello, World!"
	if len(args) > 2 || len(args) < 1 {
		return &object.String{Value: "not enough arguments"}
	}

	return &object.String{Value: value}
}

// Register the Hello function as an object.Extended
// This allows the function to be called from the Monkey interpreter
func Register() {
	object.RegisterFunction("hello", object.Extended{Fn: Hello})
}

// Build the plugin
// go build -buildmode=plugin -o extensions/hello.so hello.go
