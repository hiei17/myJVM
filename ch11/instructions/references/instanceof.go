package references

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"
	"jvmgo/ch11/rtda/heap"
)

// Determine if object is of given type
type INSTANCE_OF struct{ base.Index16Instruction }//一个操作数

func (self *INSTANCE_OF) Execute(frame *rtda.Frame) {

	//mark 1 弹出栈顶对象
	stack := frame.OperandStack()
	topObject := stack.PopRef()
	//如果是null，则把0推入操作数栈
	if topObject == nil {
		stack.PushInt(0)
		return
	}

	//mark 2 拿到操作数指定的类
	cp := frame.Method().Class().ConstantPool()
	classRef := cp.GetConstant(self.Index).(*heap.ClassRef)
	class := classRef.ResolvedClass()

	//mark 3 类型判断
	if topObject.IsInstanceOf(class) {
		stack.PushInt(1)
	} else {
		stack.PushInt(0)
	}
}
