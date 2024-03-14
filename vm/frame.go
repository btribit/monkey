// vm/frame.go

package vm

import (
	"monkey/code"
	"monkey/object"
)

// Frame is a struct to hold the information about the frame
type Frame struct {
	fn          *object.CompiledFunction
	ip          int
	basePointer int
}

// NewFrame is a function to create a new frame
func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	f := &Frame{
		fn:          fn,
		ip:          -1,
		basePointer: basePointer,
	}
	return f
}

// Instructions is a function to return the instructions
func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
