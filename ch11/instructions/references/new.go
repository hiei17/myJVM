package references

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"

	"jvmgo/ch11/rtda/heap"
)

/*

new指令的操作数是一个uint16索引，来自字节码。
通过这个索引，可以从当前类的运行时常量池中找到一个类符号引用。
解析这个类符号引用，拿到类数据，然后创建对象，并把对象引用推入栈顶，
new指令的工作就完成了。

*/


// Create new object//操作数来自字节码 是pool的坐标 里面是类名
type NEW struct{ base.Index16Instruction }



func (self *NEW) Execute(frame *rtda.Frame) {

	//mark 1.得操作数指定的类
	///拿到本类的常量池
	thisClassPool := frame.Method().Class().ConstantPool()

	//self.Index是操作数
	// 通过它从当前类的运行时常量池中找到一个类符号引用
	classRef := thisClassPool.GetConstant(self.Index).(*heap.ClassRef)//强转
	//解析这个类符号引用，拿到类数据
	class := classRef.ResolvedClass()//得引用的类

	//先判 断类的初始化是否已经开始，
	// 如果还没有，则需要调用类的初始化方法，并终止指令执行。
	// 但是由于此时指令已经执行到了一半，也就是说当前帧的nextPC字段已经指向下一条指令，
	// 所以需要修改 nextPC，让它重新指向当前指令。
	if !class.InitStarted() {
		frame.RevertNextPC()

		base.InitClass(frame.Thread(), class)
		return
	}


	/*因为接口和抽象类都不能实例化，所以如果解析后的类是接
	口或抽象类，按照Java虚拟机规范规定，需要抛出InstantiationError
	异常。*/
	if class.IsInterface() || class.IsAbstract() {
		panic("java.lang.InstantiationError")
	}

	//mark 2 new出来一个实例 压入栈
	//实例 里面有:对类的引用class 和new好一个实例变量数组(还没赋值) field
	//哪个实例变量在field数组里面是哪个 在class的field里面有个slotId字段存了(class加载后link里面就算好了
	oneObjectRef := class.NewObject()
	frame.OperandStack().PushRef(oneObjectRef)//入操作数栈

}
