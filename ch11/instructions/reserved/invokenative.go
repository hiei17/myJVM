package reserved

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/native"
import _ "jvmgo/ch11/native/java/lang"
import _ "jvmgo/ch11/native/sun/misc"
/*如果没有任何包依赖lang包，它就不会被编译进可执行文件，
上面的本地方法也就不会被注册。所以需要一个地方导入lang包，
把它放在invokenative.go文件中。由于没有显示使用lang中的变量或
函数，所以必须在包名前面加上下划线，否则无法通过编译。这个
技术在Go语言中叫作“import for side effect”。*/

//0xFE指令

// Invoke native method
type INVOKE_NATIVE struct{ base.NoOperandsInstruction }

func (self *INVOKE_NATIVE) Execute(frame *rtda.Frame) {

	method := frame.Method()
	className := method.Class().Name()
	methodName := method.Name()
	methodDescriptor := method.Descriptor()

//根据类名、方法名和方法描述符从本地方法注册表中查找本地方法实现
	nativeMethod := native.FindNativeMethod(className, methodName, methodDescriptor)

	methodInfo := className + "." + methodName + methodDescriptor
	if nativeMethod == nil {
		panic("java.lang.UnsatisfiedLinkError: " + methodInfo)
	}

	nativeMethod(frame)

}

