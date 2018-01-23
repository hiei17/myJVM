package references

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"

	"jvmgo/ch11/rtda/heap"
)

//接口方法 就是类被声明为接口 调用接口方法  ((Runnable) demo).run();
//动态绑定  需要看this传参确定实际的类型 找到实际要执行的方法


// Invoke interface method
type INVOKE_INTERFACE struct {

	index uint
	// count uint8
	// zero uint8
}

//拿到操作数 其实只用第一个数
func (self *INVOKE_INTERFACE) FetchOperands(reader *base.BytecodeReader) {

	//前两字节的含义和其他指令相同，是个uint16运行时常量池索引。
	self.index = uint(reader.ReadUint16())

	//给方法传递参数需要的slot数，
	// 其含义和给Method结构体定义的argSlotCount字段相同。
	// 这个数是可以根据方法描述符计算出来的，它的存在仅仅是因为历史原因。
	reader.ReadUint8() // count

	//是留给Oracle的某些Java虚拟机实现用的，它的值必须是0。
	// 该字节的存在是为了保证Java虚拟机可以向后兼容。
	reader.ReadUint8() // must be 0
}

//动态绑定
func (self *INVOKE_INTERFACE) Execute(frame *rtda.Frame) {

	cp := frame.Method().Class().ConstantPool()
	methodRef := cp.GetConstant(self.index).(*heap.InterfaceMethodRef)
	resolvedMethod := methodRef.ResolvedInterfaceMethod()

	//不能是静态 或者 是私有方法
	if resolvedMethod.IsStatic() || resolvedMethod.IsPrivate() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	thisObject := frame.OperandStack().GetRefFromTop(resolvedMethod.ArgSlotCount() - 1)

	if thisObject == nil {
		panic("java.lang.NullPointerException") // todo
	}

	//mark 引用所指对象的类 没有实现 解析出来的接口
	thisClass := thisObject.Class()
	if !thisClass.IsImplements(methodRef.ResolvedClass()) {
		panic("java.lang.IncompatibleClassChangeError")
	}

	//mark  最终要调用的方法
	methodToBeInvoked := heap.LookupMethodInClass(thisClass,
		methodRef.Name(), methodRef.Descriptor())

	if methodToBeInvoked == nil || methodToBeInvoked.IsAbstract() {
		panic("java.lang.AbstractMethodError")
	}

	//不是public
	if !methodToBeInvoked.IsPublic() {
		panic("java.lang.IllegalAccessError")
	}

	//一切正常，调用方法
	base.InvokeMethod(frame, methodToBeInvoked)
}



