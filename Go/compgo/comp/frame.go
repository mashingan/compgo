package comp

type Frame struct {
	cl          *Closure
	ip          int
	basePointer int
}

func NewFrame(cl *Closure, basePointer int) *Frame {
	return &Frame{cl, 0, basePointer}
}

func (f *Frame) Instructions() Instructions {
	return f.cl.Fn.Instructions
}

func (f *Frame) SetInstructions(ins Instructions) {
	f.cl.Fn.Instructions = ins
}
