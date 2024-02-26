// compiler/compiler.go

package compiler

import (
	"fmt"
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
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	}

	return nil
}

// addConstant adds a constant to the compiler's constant pool and returns its position
func (c *compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// emit appends the given instructions to the compiler's instruction stream
func (c *compiler) emit(op code.Opcode, operands ...int) int {
	instruction := code.Make(op, operands...)
	def, err := code.Lookup(instruction[0])
	if err != nil {
		return c.addInstruction(instruction)
	}
	fmt.Printf("emit %s\n", def.Name)
	position := c.addInstruction(instruction)
	return position
}

// addInstruction adds the given instructions to the compiler's instruction stream
func (c *compiler) addInstruction(ins []byte) int {
	positionNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return positionNewInstruction
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
