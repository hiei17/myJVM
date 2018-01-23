package rtda

import "jvmgo/ch11/rtda/heap"



//栈帧
type Frame struct {
	lower *Frame//用来实现链表数据结构
	localVars LocalVars//局部变量表
	operandStack *OperandStack//操作数栈
	method *heap.Method//产生本Frame的方法

	//为了实现跳转指令而添加的
	thread *Thread//所在的线程
	nextPC int
}
//局部变量表大小 操作数栈深度是由 编译器预先计算好的，存储在method_info结构的Code属性中
func newFrame(thread *Thread, method *heap.Method) *Frame {
	return &Frame{
		thread: thread,
		method: method,
		localVars: newLocalVars(method.MaxLocals()),
		operandStack: newOperandStack(method.MaxStack()),
	}
}

// getters & setters
func (self *Frame) LocalVars() LocalVars {
	return self.localVars
}

//操作数栈
func (self *Frame) OperandStack() *OperandStack {
	return self.operandStack
}

func (self *Frame) Thread() *Thread {
	return self.thread
}
func (self *Frame) NextPC() int {
	return self.nextPC
}
func (self *Frame) SetNextPC(nextPC int) {
	self.nextPC = nextPC
}


//此时指令已经执行到了一半，也 就是说当前帧的nextPC字段已经指向下一条指令，所以需要修改 nextPC，让它重新指向当前指令
func (self *Frame) RevertNextPC() {
	self.nextPC = self.thread.pc
}

func (self *Frame) Method() *heap.Method {
	return self.method
}

func (self *Frame) GetFieldByConstantIndex(index uint) *heap.Field{

	currentClass := self.method.Class()//每个方法里面都有类
	classPool := currentClass.ConstantPool()//类的运行时常量池
	//字段引用
	fieldRef := classPool.GetConstant(index).(*heap.FieldRef)
	//得字段对象 有存着字段所在类 描述 名字   得字段对象
	return fieldRef.ResolvedField()

}
