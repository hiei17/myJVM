package control

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"
)



/*
tableswitch
<0-3 byte pad>
defaultbyte1
defaultbyte2
defaultbyte3
defaultbyte4
lowbyte1
lowbyte2
lowbyte3
lowbyte4
highbyte1
highbyte2
highbyte3
highbyte4
jump offsets...
*/
// Access jump table by index and jump
type TABLE_SWITCH struct {
	defaultOffset int32
	low           int32
	high          int32
	jumpOffsets   []int32
}


func (self *TABLE_SWITCH) FetchOperands(reader *base.BytecodeReader) {
	//tableswitch指令操作码的后面有0~3字节的padding，
	// 以保证 defaultOffset在字节码中的地址是4的倍数
	reader.SkipPadding()//跳过几个

	//默认情况下执行跳转所需的字节码偏移量
	self.defaultOffset = reader.ReadInt32()//default

	//记录case的取值范围
	self.low = reader.ReadInt32()
	self.high = reader.ReadInt32()

	//一个索引表，里面存放high-low+1个int值，对应各种case情况下，执行跳转所需的字节码偏移量
	jumpOffsetsCount := self.high - self.low + 1
	self.jumpOffsets = reader.ReadInt32s(jumpOffsetsCount)
}
/*

栈中弹出一个int变量
是否在low和high给定的范围之内。
如果在，则从jumpOffsets表中查出偏移量进行跳转，
否则按照defaultOffset跳转按照defaultOffset跳转
*/

func (self *TABLE_SWITCH) Execute(frame *rtda.Frame) {
	index := frame.OperandStack().PopInt()//switch(index)

	var offset int
	if index >= self.low && index <= self.high {
		//在case范围内 去跳转表查偏移量
		offset = int(self.jumpOffsets[index-self.low])
	} else {
		//否则按照defaultOffset跳转按照defaultOffset跳转
		offset = int(self.defaultOffset)
	}

	base.Branch(frame, offset)
}
