package control

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"
)


//无条件跳转


// Branch always
type GOTO struct{ base.BranchInstruction }

func (self *GOTO) Execute(frame *rtda.Frame) {
	base.Branch(frame, self.Offset)
}
