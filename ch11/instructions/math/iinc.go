package math

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"
)

//给局部变量表中的int变量增加常量值，
// 局部变量表索引和常量值 都由 指令的操作数提供

// Increment local variable by constant
type IINC struct {
	Index uint
	Const int32
}

func (self *IINC) FetchOperands(reader *base.BytecodeReader) {
	self.Index = uint(reader.ReadUint8())
	self.Const = int32(reader.ReadInt8())
}

func (self *IINC) Execute(frame *rtda.Frame) {
	localVars := frame.LocalVars()
	val := localVars.GetInt(self.Index)
	val += self.Const
	localVars.SetInt(self.Index, val)
}
