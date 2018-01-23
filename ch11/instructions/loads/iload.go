package loads

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"
)
//从 局部变量表 获取变量，推入 操作数栈顶
//公共函数
func _iload(frame *rtda.Frame, index uint) {
	val := frame.LocalVars().GetInt(index)
	frame.OperandStack().PushInt(val)
}

// Load int from local variable
type ILOAD struct{ base.Index8Instruction }//局部变量表索引 在操作数中
func (self *ILOAD) Execute(frame *rtda.Frame) {
	_iload(frame, uint(self.Index))
}

//++++++++++++++++++++++++++++++++++++以下是把局部变量表第0 1 2 3 slot 入操作数栈顶

type ILOAD_0 struct{ base.NoOperandsInstruction }
func (self *ILOAD_0) Execute(frame *rtda.Frame) {
	_iload(frame, 0)
}

type ILOAD_1 struct{ base.NoOperandsInstruction }
func (self *ILOAD_1) Execute(frame *rtda.Frame) {
	_iload(frame, 1)
}

type ILOAD_2 struct{ base.NoOperandsInstruction }
func (self *ILOAD_2) Execute(frame *rtda.Frame) {
	_iload(frame, 2)
}

type ILOAD_3 struct{ base.NoOperandsInstruction }
func (self *ILOAD_3) Execute(frame *rtda.Frame) {
	_iload(frame, 3)
}


