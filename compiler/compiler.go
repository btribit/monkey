// compiler/compiler.go

package compiler

import (
	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

type compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

func New() *compiler {
	return &compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *compiler) Compile(node ast.Node) error {
	return nil
}

func (c *compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
