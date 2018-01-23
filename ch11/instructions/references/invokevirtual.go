package references

import "fmt"
import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/rtda/heap"

// Invoke instance method; dispatch based on class
type INVOKE_VIRTUAL struct{ base.Index16Instruction }

func (self *INVOKE_VIRTUAL) Execute(frame *rtda.Frame) {

	//找到表面上要调用的方法
	currentClass := frame.Method().Class()
	cp := currentClass.ConstantPool()
	methodRef := cp.GetConstant(self.Index).(*heap.MethodRef)
	resolvedMethod := methodRef.ResolvedMethod()

	if resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	thisObjectPosition := resolvedMethod.ArgSlotCount()-1
	stack := frame.OperandStack()
	thisObject := stack.GetRefFromTop(thisObjectPosition)

	if thisObject == nil {
		// hack!
		if methodRef.Name() == "println" {
			_println(stack, methodRef.Descriptor())
			return
		}

		panic("java.lang.NullPointerException")
	}

	//访问限制检查
	if resolvedMethod.IsProtected() &&
	resolvedMethod.Class().IsSuperClassOf(currentClass) &&
	resolvedMethod.Class().GetPackageName() != currentClass.GetPackageName() &&
	thisObject.Class() != currentClass &&
	!thisObject.Class().IsSubClassOf(currentClass) {

		panic("java.lang.IllegalAccessError")
	}

	//mark 类用传参this的类 找到实际上调用的方法
	methodToBeInvoked := heap.LookupMethodInClass(thisObject.Class(), methodRef.Name(), methodRef.Descriptor())

	if methodToBeInvoked == nil || methodToBeInvoked.IsAbstract() {
		panic("java.lang.AbstractMethodError")
	}

	//调用
	base.InvokeMethod(frame, methodToBeInvoked)
}

//传参:操作数栈,被打印的数据的描述符
// hack!
func _println(stack *rtda.OperandStack, descriptor string) {

	//不同类型的数据打印
	switch descriptor {
		case "(Z)V":
			fmt.Printf("%v\n", stack.PopInt() != 0)
		case "(C)V":
			fmt.Printf("%c\n", stack.PopInt())
		case "(I)V", "(B)V", "(S)V":
			fmt.Printf("%v\n", stack.PopInt())
		case "(F)V":
			fmt.Printf("%v\n", stack.PopFloat())
		case "(J)V":
			fmt.Printf("%v\n", stack.PopLong())
		case "(D)V":
			fmt.Printf("%v\n", stack.PopDouble())
		//mark 字符串打印
		case "(Ljava/lang/String;)V"://打印字符串
			stringObject := stack.PopRef()
			goStr := heap.GoString(stringObject)//从字符串对象里面得到字符串
			fmt.Println(goStr)

		default:
			panic("println: " + descriptor)
		}
	stack.PopRef()
}
