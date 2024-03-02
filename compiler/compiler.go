// compiler/compiler.go

package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type compiler struct {
	instructions code.Instructions
	constants    []object.Object

	lastInstruction EmittedInstruction
	prevInstruction EmittedInstruction
}

func New() *compiler {
	return &compiler{
		instructions:    code.Instructions{},
		constants:       []object.Object{},
		lastInstruction: EmittedInstruction{},
		prevInstruction: EmittedInstruction{},
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

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// emit an OpJumpNotTruthy with a bogus value
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		jumpPos := c.emit(code.OpJump, 9999)

		afterConsequencePos := len(c.instructions)
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}
		}

		afterAlternativePos := len(c.instructions)
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.BlockStatement:
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

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}

			c.emit(code.OpGreaterThan)
			return nil
		}

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
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	}

	return nil
}

// replaceInstruction replaces the instruction at the given position with the given instructions
func (c *compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i]
	}
}

// changeOperand changes the operand of the instruction at the given position
func (c *compiler) changeOperand(opPos int, operand int) {
	opcode := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(opcode, operand)
	c.replaceInstruction(opPos, newInstruction)
}

// lastInstructionIsPop returns true if the last instruction is a pop
func (c *compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

// removeLastPop removes the last pop instruction
func (c *compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.prevInstruction
}

// addConstant adds a constant to the compiler's constant pool and returns its position
func (c *compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// emit appends the given instructions to the compiler's instruction stream
func (c *compiler) emit(op code.Opcode, operands ...int) int {
	instruction := code.Make(op, operands...)
	position := c.addInstruction(instruction)

	def, err := code.Lookup(instruction[0])
	if err != nil {
		return c.addInstruction(instruction)
	}
	fmt.Printf("emit %s\n", def.Name)

	c.setLastInstruction(op, position)
	return position
}

// setLastInstruction sets the last instruction and the previous instruction
func (c *compiler) setLastInstruction(op code.Opcode, position int) {
	c.prevInstruction = c.lastInstruction
	c.lastInstruction = EmittedInstruction{Opcode: op, Position: position}
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
