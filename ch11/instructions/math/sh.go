package math

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"
)

//位移指令

// Shift left int
type ISHL struct{ base.NoOperandsInstruction }

func (self *ISHL) Execute(frame *rtda.Frame) {

	stack := frame.OperandStack()
	v2 := stack.PopInt()//要进行位移操作的变量
	v1 := stack.PopInt()//要移位多少比特
	//Go语言位移操作符右侧必须是无符号整数，所以需要对v2进行类型转换
	s := uint32(v2) & 0x1f//int变量只有32位，所以只取v2的前5个比特就足够表示位移位数了
	result := v1 << s
	stack.PushInt(result)//把结果推入操作数栈
}

// Arithmetic shift right int
type ISHR struct{ base.NoOperandsInstruction }

func (self *ISHR) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	s := uint32(v2) & 0x1f
	result := v1 >> s
	stack.PushInt(result)
}

// Logical shift right int
type IUSHR struct{ base.NoOperandsInstruction }

func (self *IUSHR) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	s := uint32(v2) & 0x1f
	result := int32(uint32(v1) >> s)
	stack.PushInt(result)
}

// Shift left long
type LSHL struct{ base.NoOperandsInstruction }

func (self *LSHL) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopLong()
	s := uint32(v2) & 0x3f
	result := v1 << s
	stack.PushLong(result)
}

// Arithmetic shift right long
type LSHR struct{ base.NoOperandsInstruction }
func (self *LSHR) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopLong()
	s := uint32(v2) & 0x3f
	result := v1 >> s
	stack.PushLong(result)
}

// Logical shift right long
type LUSHR struct{ base.NoOperandsInstruction }

func (self *LUSHR) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopLong()
	s := uint32(v2) & 0x3f
	result := int64(uint64(v1) >> s)
	stack.PushLong(result)
}
