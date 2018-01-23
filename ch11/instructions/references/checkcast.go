
package references
import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/rtda/heap"
// Check whether object is of given type
type CHECK_CAST struct{ base.Index16Instruction }
/*

checkcast指令和instanceof指令很像，
区别在于：
checkcast则不改变操作数栈（如果判断失败，直接抛出ClassCastException异常）。
c
*/

func (self *CHECK_CAST) Execute(frame *rtda.Frame) {

	//mark 1. 得栈顶的实例
	topObject := peek(frame  )
	//null都可以 返回
	if topObject == nil {
		return
	}

	//mark 2 .得操作数 指明的类
	class :=getClass(self.Index,frame)

	//mark 3 .要是栈顶的对象 不是指定那个类的  就强转异常
	if !topObject.IsInstanceOf(class) {
		panic("java.lang.ClassCastException")
	}
}
func getClass(i uint,frame *rtda.Frame) *heap.Class {
	cp := frame.Method().Class().ConstantPool()
	classRef := cp.GetConstant(i).(*heap.ClassRef)
	return classRef.ResolvedClass()
}
func peek(frame *rtda.Frame)  *heap.Object  {
	stack := frame.OperandStack()
	//先从操作数栈中弹出对象引用，再推回去，这样就不会改变操作数栈的状态。
	ref := stack.PopRef()
	stack.PushRef(ref)
	return ref
}

