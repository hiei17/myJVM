package references

import "reflect"
import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/rtda/heap"
// Throw exception or error
type ATHROW struct{ base.NoOperandsInstruction }

func (self *ATHROW) Execute(frame *rtda.Frame) {

	//异常对象引用
	ex := frame.OperandStack().PopRef()
	if ex == nil {
		panic("java.lang.NullPointerException")
	}
	thread := frame.Thread()

	//可以找到并跳转到异常处理代码
	if !findAndGotoExceptionHandler(thread, ex) {

		//遍历完Java虚拟机栈还是找不到catch
		//打印
		handleUncaughtException(thread, ex)
	}
}


//找到并跳转到异常处 理代码
func findAndGotoExceptionHandler(thread *rtda.Thread, ex *heap.Object) bool {

	//从当前帧开始，遍历Java虚拟机栈，查找方法的异常处理表。
	for {
		//当前帧的所属方法 有没有写过catch 这个位置的这个类型异常
		frame := thread.CurrentFrame()
		pc := frame.NextPC() - 1
		handlerPC := frame.Method().FindExceptionHandler(ex.Class(), pc)

		//如果找到了异常处理项
		if handlerPC > 0 {
			stack := frame.OperandStack()
			//操作数栈清空
			stack.Clear()

			//异常对象引用推入栈顶
			stack.PushRef(ex)

			//跳转到异常处理代码
			frame.SetNextPC(handlerPC)
			return true
		}

		//找不到异常处理项，则把frame弹出，继续遍历。
		thread.PopFrame()
		if thread.IsStackEmpty() {
			break
		}
	}
	return false
}

//打印出Java虚拟机栈信息
func handleUncaughtException(thread *rtda.Thread, ex *heap.Object) {
	//Java虚拟机栈清空  Java虚拟机栈已经空了，所以解释器也就终止执行了
	thread.ClearStack()

	//得异常对象的这个字段值 转成go字符串
	jMsg := ex.GetRefVar("detailMessage", "Ljava/lang/String;")
	goMsg := heap.GoString(jMsg)
	//打印 异常类名:异常信息
	println(ex.Class().JavaName() + ": " + goMsg)

	//Java虚拟机栈信息 被放在extra里面了
	stes := reflect.ValueOf(ex.Extra())

	for i := 0; i < stes.Len(); i++ {

		ste := stes.Index(i).Interface().(interface {
			String() string
		})
		println("\tat " + ste.String())
	}
}

