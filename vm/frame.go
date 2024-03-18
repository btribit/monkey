// vm/frame.go

package vm

import (
	"monkey/code"
	"monkey/object"
)

// Frame is a struct to hold the information about the frame
type Frame struct {
	cl          *object.Closure
	ip          int
	basePointer int
}

// NewFrame is a function to create a new frame
func NewFrame(cl *object.Closure, basePointer int) *Frame {
	f := &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,
	}
	return f
}

// Instructions is a function to return the instructions
func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
