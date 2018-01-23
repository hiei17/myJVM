package base

import (
	"jvmgo/ch11/rtda"
	"jvmgo/ch11/rtda/heap"
)
//方法执行
/*

Java虚拟机要给这个方法创建一个新的帧并把它推入Java虚拟机栈顶，
传递参数。

*/
//各个方法执行指令的公共部分
func InvokeMethod(invokerFrame *rtda.Frame, method *heap.Method) {

	thread := invokerFrame.Thread()
	newFrame := thread.NewFrame(method)
	thread.PushFrame(newFrame)

	//传递参数
	//方法的参数在局部 变量表中占用多少位置
	argSlotCount:= int(method.ArgSlotCount())//method对象从classfile产生时就就根据描述符和是否static算好了
	/*

		注意，这个数量并不一定等于从Java代码
		中看到的参数个数，原因有两个。
		第一，long和double类型的参数要
		占用两个位置。
		第二，对于实例方法，Java编译器会在参数列表的
	*/

	argSlotSlot := argSlotCount

	//按顺序从老帧的操作数栈 到 新帧的局部变量表  要是实例方法 this在新帧局部变量表第一个位置
	if argSlotSlot > 0 {
		for i := argSlotSlot - 1; i >= 0; i-- {

			//依次把这n个变量从调用者的操作数栈中弹出
			slot := invokerFrame.OperandStack().PopSlot()

			//放进被调用方法的局部变量表中
			newFrame.LocalVars().SetSlot(uint(i), slot)
		}
	}

	//Java虚拟机规范并没有规定如何实现和调用本地方法
	if method.IsNative() {
		//todo
	}

}
