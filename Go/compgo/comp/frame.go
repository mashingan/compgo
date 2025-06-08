package comp

type Frame struct {
	fn *CompiledFunction
	ip int
}

func NewFrame(fn *CompiledFunction) *Frame {
	return &Frame{fn, 0}
}

func (f *Frame) Instructions() Instructions {
	return f.fn.Instructions
}

func (f *Frame) SetInstructions(ins Instructions) {
	f.fn.Instructions = ins
}
