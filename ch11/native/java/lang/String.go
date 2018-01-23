package lang

import "jvmgo/ch11/native"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/rtda/heap"

//字符串池

func init() {
	native.Register("java/lang/String", "intern", "()Ljava/lang/String;", intern)
}

// public native String intern();
func intern(frame *rtda.Frame) {
	this := frame.LocalVars().GetThis()
	interned := heap.InternString(this)
	frame.OperandStack().PushRef(interned)
}


