package extended

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"
)



//和基本goto的 不同是 ReadInt32
// Branch always (wide index)
type GOTO_W struct {
	offset int
}

func (self *GOTO_W) FetchOperands(reader *base.BytecodeReader) {
	self.offset = int(reader.ReadInt32())//TODO
}
func (self *GOTO_W) Execute(frame *rtda.Frame) {
	base.Branch(frame, self.offset)
}
