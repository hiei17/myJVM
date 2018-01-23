package references

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"

//int x = a1.length; 获取数组长度 压入栈
// Get length of array
type ARRAY_LENGTH struct{ base.NoOperandsInstruction }

func (self *ARRAY_LENGTH) Execute(frame *rtda.Frame) {
	//从操作数栈顶弹出的数组
	stack := frame.OperandStack()
	arrObject := stack.PopRef()
	if arrObject == nil {
		panic("java.lang.NullPointerException")
	}

	arrLen := arrObject.ArrayLength()
	stack.PushInt(arrLen)
}
