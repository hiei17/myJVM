package lang

import "jvmgo/ch11/native"
import "jvmgo/ch11/rtda"
import (
	"jvmgo/ch11/rtda/heap"

)

func init() {
	native.Register("java/lang/Throwable", "fillInStackTrace",
		"(I)Ljava/lang/Throwable;", fillInStackTrace)
}


//记录Java虚拟机栈帧信息 正在执行的信息
type StackTraceElement struct {
	fileName string
	className string
	methodName string
	lineNumber int//正在执行哪行代码
}


// private native Throwable fillInStackTrace(int dummy);
func fillInStackTrace(frame *rtda.Frame) {

	this := frame.LocalVars().GetThis()
	frame.OperandStack().PushRef(this)

	stes := createStackTraceElements(this, frame.Thread())
	this.SetExtra(stes)
}

func createStackTraceElements(tObj *heap.Object, thread *rtda.Thread)[]*StackTraceElement {

	//由于栈顶两帧正在执行fillInStackTrace（int）和fillInStackTrace（）方法，所以需要跳过这两帧。
	// 下面的几帧正在执行异常类的构造函数，所以也要跳过，具体要跳过多少帧数则要看异常类的继承层次
	skip := distanceToObject(tObj.Class()) + 2

	/*计算好需要跳过的帧之后，调用Thread结构体的GetFrames（）
方法拿到完整的Java虚拟机栈，然后reslice一下就是真正需要的帧。
GetFrames（）方法只是调用了Stack结构体的getFrames（）方法，代码
如下：*/
	frames := thread.GetFrames()[skip:]

	stes := make([]*StackTraceElement, len(frames))
	for i, frame := range frames {
		stes[i] = createStackTraceElement(frame)
	}
	return stes
}

func distanceToObject(class *heap.Class) int {
	distance := 0
	for c := class.SuperClass(); c != nil; c = c.SuperClass() {
		distance++
	}
	return distance
}


//根据帧 创建StackTraceElement实例
func createStackTraceElement(frame *rtda.Frame) *StackTraceElement {

	method := frame.Method()
	class := method.Class()

	return &StackTraceElement{
		fileName: class.SourceFile(),
		className: class.JavaName(),
		methodName: method.Name(),
		lineNumber: method.GetLineNumber(frame.NextPC() - 1),
	}
}
