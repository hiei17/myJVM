package base

import "jvmgo/ch11/rtda"
import "jvmgo/ch11/rtda/heap"

// jvms 5.5
func InitClass(thread *rtda.Thread, class *heap.Class) {

	//类的initStarted状态设置成true以免进入死循环
	class.StartInit()

	//准备 执行类的初始化方法
	scheduleClinit(thread, class)

	//要先执行超类的初始化方法
	initSuperClass(thread, class)
}

func scheduleClinit(thread *rtda.Thread, class *heap.Class) {

	//找到方法名为<clinit> 描述符为()V 的方法
	clinit := class.GetClinitMethod()
	if clinit != nil {
		// exec <clinit>
		newFrame := thread.NewFrame(clinit)
		thread.PushFrame(newFrame)
	}
}

func initSuperClass(thread *rtda.Thread, class *heap.Class) {
	if !class.IsInterface() {
		superClass := class.SuperClass()
		if superClass != nil && !superClass.InitStarted() {
			InitClass(thread, superClass)
		}
	}
}
