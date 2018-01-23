package constants

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"
)


//啥也不做的指令
// Do nothing//复用base.NoOperandsInstruction 能用base.NoOperandsInstruction 的一切:	相当于继承
type NOP struct{ base.NoOperandsInstruction }

func (self *NOP) Execute(frame *rtda.Frame) {
	// really do nothing
}


//base.NoOperandsInstruction 里面已经有FetchOperands函数 现在又写了Execute
//Instruction定义的空函数都实现了  就是Instruction了
