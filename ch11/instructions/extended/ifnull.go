package extended

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"
)

//根据引用是否是null进行跳转


// Branch if reference is null
type IFNULL struct{ base.BranchInstruction }
func (self *IFNULL) Execute(frame *rtda.Frame) {
	ref := frame.OperandStack().PopRef()
	if ref == nil {
		base.Branch(frame, self.Offset)
	}
}

// Branch if reference not null
type IFNONNULL struct{ base.BranchInstruction }
func (self *IFNONNULL) Execute(frame *rtda.Frame) {
	ref := frame.OperandStack().PopRef()
	if ref != nil {
		base.Branch(frame, self.Offset)
	}
}
