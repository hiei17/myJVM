package native
import "jvmgo/ch11/rtda"

//注册和查找本地方法

//本地方法
//这个frame参数就是本地方法的工作空间 也就是连接Java虚拟机和Java类库的桥梁
type NativeMethod func(frame *rtda.Frame)

//registry变量是个哈希表，值是具体的本地方法实现
var registry = map[string]NativeMethod{}

//加到registry
func Register(className, methodName, methodDescriptor string, nativeMethod NativeMethod) {
	key := className + "~" + methodName + "~" + methodDescriptor//类名、方法名和方法描述符加在一起才能唯一确定一个方法
	registry[key] = nativeMethod
}

//来key  返回 本地方法
func FindNativeMethod(className, methodName, methodDescriptor string) NativeMethod {
	key := className + "~" + methodName + "~" + methodDescriptor
	if method, ok := registry[key]; ok {
		return method
	}
	if methodDescriptor == "()V" && methodName == "registerNatives" {
		return emptyNativeMethod
	}
	return nil
}

/*jva.lang.Object等类是通过一个叫作registerNatives（）的本地方法来
注册其他本地方法的。在本章和后面的章节中，将自己注册所有的
本地方法实现。所以像registerNatives（）这样的方法就没有太大的用
处。为了避免重复代码，这里统一处理，如果遇到这样的本地方
法，就返回一个空的实现，代码如下：*/
func emptyNativeMethod(frame *rtda.Frame) {
	// do nothing
}

