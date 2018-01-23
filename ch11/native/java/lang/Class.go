package lang

import "jvmgo/ch11/native"
import "jvmgo/ch11/rtda"
import (
	"jvmgo/ch11/rtda/heap"

)
func init() {
	native.Register(
		"java/lang/Class",
		"getPrimitiveClass",
		"(Ljava/lang/String;)Ljava/lang/Class;",
		getPrimitiveClass)
	native.Register(
		"java/lang/Class",
		"getName0",
		"()Ljava/lang/String;",
		getName0)
	native.Register("java/lang/Class",
		"desiredAssertionStatus0",
		"(Ljava/lang/Class;)Z",
		desiredAssertionStatus0)
}

//基本类型的包装类在初 始化时会调用
//如Integer类里面初始化会执行:
//public static final Class<Integer> TYPE = (Class<Integer>) Class.getPrimitiveClass("int");
//让类加载器加载Integer基本类 得其class对象 入栈
func getPrimitiveClass(frame *rtda.Frame) {
	nameObj := frame.LocalVars().GetRef(0)//"int"java的String对象
	name := heap.GoString(nameObj)//转成Go字符串

	//基本类型的类已经加载到了方法区中，直接调用类加载器的LoadClass（）方法获取即可。
	loader := frame.Method().Class().Loader()
	class := loader.LoadClass(name).JClass()
	frame.OperandStack().PushRef(class)
}

//Class.getName（）方法 内部是调用这个的
//  Class c=Integer.class;
//c.getName();//内部调用getName0
func getName0(frame *rtda.Frame) {
	thisObject := frame.LocalVars().GetThis()
	class := thisObject.Extra().(*heap.Class)
	name := class.JavaName()//如java.lang.Object
	nameObj := heap.JString(class.Loader(), name)//转成Java字符串
	frame.OperandStack().PushRef(nameObj)
}


//Character类是基本类型char的包装类，
// 在初始化时会调用Class.desiredAssertionStatus0（）
//hack!
// private static native boolean desiredAssertionStatus0(Class<?> clazz);
func desiredAssertionStatus0(frame *rtda.Frame) {
	frame.OperandStack().PushBoolean(false)
}
