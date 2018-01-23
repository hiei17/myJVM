package constants

import "jvmgo/ch11/instructions/base"
import (
	"jvmgo/ch11/rtda"
	"jvmgo/ch11/rtda/heap"
)
/*

从运行时常量池中加载常量值，并把它推入栈。
ldc: int、float和字符串常量
ldc_w:java.lang.Class实例或者MethodType和MethodHandle实例。
ldc2_w: long和double常量。
ldc和ldc_w指令的区别仅在于操作数的宽度。
本章只处理int、float、long和double常量。

*/

type LDC struct{ base.Index8Instruction }
type LDC_W struct{ base.Index16Instruction }
type LDC2_W struct{ base.Index16Instruction }

func (self *LDC) Execute(frame *rtda.Frame) {
	_ldc(frame, self.Index)
}
func (self *LDC_W) Execute(frame *rtda.Frame) {
	_ldc(frame, self.Index)
}

//公共的
func _ldc(frame *rtda.Frame, index uint) {
	class := frame.Method().Class()
	classPool := class.ConstantPool()
	constant := classPool.GetConstant(index)
	stack := frame.OperandStack()
	switch constant.(type) {
		case int32: stack.PushInt(constant.(int32))
		case float32: stack.PushFloat(constant.(float32))
		case string:
			//转成Java字符串实例
			internedStr := heap.JString(class.Loader(), constant.(string))
			stack.PushRef(internedStr)

		//加载class对象 普通的对象需要new出来一个入栈 但是类的相应class对象就在它自己里面
		case *heap.ClassRef:
			classRef := constant.(*heap.ClassRef)
			classObj := classRef.ResolvedClass().JClass()
			stack.PushRef(classObj)
		default: panic("todo: ldc!")
	}
}

func (self *LDC2_W) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	cp := frame.Method().Class().ConstantPool()
	c := cp.GetConstant(self.Index)
	switch c.(type) {
		case int64: stack.PushLong(c.(int64))
		case float64: stack.PushDouble(c.(float64))
		//其他情况还无法处理，暂 时调用panic（）函数终止程序执行
		default: panic("java.lang.ClassFormatError")
	}
}
