package ast

import (
	"monkey/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{ // Create a new program
		Statements: []Statement{ // Initialize the statements
			&LetStatement{ // Create a new let statement
				Token: token.Token{Type: token.LET, Literal: "let"}, // Create a new token
				Name: &Identifier{ // Create a new identifier
					Token: token.Token{Type: token.IDENT, Literal: "myVar"}, // Create a new token
					Value: "myVar",                                          // Set the value
				},
				Value: &Identifier{ // Create a new identifier
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"}, // Create a new token
					Value: "anotherVar",                                          // Set the value
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" { // Check if the string representation is correct
		t.Errorf("program.String() wrong. Got %q", program.String())
	}
}
