package references

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/rtda/heap"

//调用静态方法  去常量池找到名字 拿到这个方法 调用即可

// Invoke a class (static) method
type INVOKE_STATIC struct{ base.Index16Instruction }

//类方法 编译就确定了方法
func (self *INVOKE_STATIC) Execute(frame *rtda.Frame) {

	//最简单 直接根据常量池引用找到写明的那个方法就行
	cp := frame.Method().Class().ConstantPool()
	methodRef := cp.GetConstant(self.Index).(*heap.MethodRef)
	staticMethod := methodRef.ResolvedMethod()

	if !staticMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	class := staticMethod.Class()

	//类还没有被初始化，则要先初始化该类
	if !class.InitStarted() {
		frame.RevertNextPC()
		base.InitClass(frame.Thread(), class)
		return
	}

	//找到的方法就是最终要执行的 执行就行
	base.InvokeMethod(frame, staticMethod)
}
