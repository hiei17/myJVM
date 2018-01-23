package lang

import (
	"jvmgo/ch11/native"
	"jvmgo/ch11/rtda"
)
import "unsafe"

const jlObject="java/lang/Object";
func init() {
	//注册getClass（）本地方法
	native.Register(jlObject, "getClass", "()Ljava/lang/Class;", getClass)

	native.Register(jlObject, "hashCode", "()I", hashCode)

	native.Register(jlObject, "clone", "()Ljava/lang/Object;", clone)
}

//对象.getClass（） 得其类对应的class对象 入栈
// public final native Class<?> getClass();
func getClass(frame *rtda.Frame) {
	thisObject := frame.LocalVars().GetThis()//就是调用GetRef（0）
	classObject := thisObject.Class().JClass()//类对象
	frame.OperandStack().PushRef(classObject)//入栈
}

// public native int hashCode();
func hashCode(frame *rtda.Frame) {
	this := frame.LocalVars().GetThis()
	//转换成uintptr类型，然后强制转换成int32
	hash := int32(uintptr(unsafe.Pointer(this)))
	frame.OperandStack().PushInt(hash)
}

func clone(frame *rtda.Frame) {
	this := frame.LocalVars().GetThis()
	cloneable := this.Class().Loader().LoadClass("java/lang/Cloneable")
	if !this.Class().IsImplements(cloneable) {
		panic("java.lang.CloneNotSupportedException")
	}
	frame.OperandStack().PushRef(this.Clone())
}

