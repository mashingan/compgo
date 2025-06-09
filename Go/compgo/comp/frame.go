package comp

type Frame struct {
	fn          *CompiledFunction
	ip          int
	basePointer int
}

func NewFrame(fn *CompiledFunction, basePointer int) *Frame {
	return &Frame{fn, 0, basePointer}
}

func (f *Frame) Instructions() Instructions {
	return f.fn.Instructions
}

func (f *Frame) SetInstructions(ins Instructions) {
	f.fn.Instructions = ins
}
