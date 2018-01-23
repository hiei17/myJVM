package references

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"
	"jvmgo/ch11/rtda/heap"
)

//调用无需动态绑定的实例方法:
//<init> private super

/*
因为对象是需要初始化的，所以每个类都至少有一个构造函数。
即使用户自己不定义，编译器也会自动生成一个默认构造函数。
*/

type INVOKE_SPECIAL struct{ base.Index16Instruction }

//3种情况 明明是实例方法 却不用动态绑定(看传过来的this找实际的类)  编译期间就能确定了
//1.new
//2.private方法
//3.调用了声明类的父类方法
func (self *INVOKE_SPECIAL) Execute(frame *rtda.Frame) {

	//先拿到当前类、当前常量池、方法符号引用，然后解析符号引
	//用，拿到解析后的类和方法
	currentClass := frame.Method().Class()
	cp := currentClass.ConstantPool()
	methodRef := cp.GetConstant(self.Index).(*heap.MethodRef)
	resolvedClass := methodRef.ResolvedClass()
	resolvedMethod := methodRef.ResolvedMethod()

	//构造方法<init>只有它自己能调用
	if resolvedMethod.Name() == "<init>" &&//是 调用的构造方法
		resolvedMethod.Class() != resolvedClass {//调用类不是方法所在类
		panic("java.lang.NoSuchMethodError")
	}

	if resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	//从操作数栈中弹出this引用  GetRefFromTop返回距离操作数栈顶n个单元格的引用变量
	thisObject := frame.OperandStack().GetRefFromTop(resolvedMethod.ArgSlotCount() - 1)
	if thisObject == nil {
		panic("java.lang.NullPointerException")
	}

	//访问限制  Protect
	if resolvedMethod.IsProtected() &&
		resolvedMethod.Class().IsSuperClassOf(currentClass) &&
		resolvedMethod.Class().GetPackageName() != currentClass.GetPackageName() &&
		thisObject.Class() != currentClass &&
		!thisObject.Class().IsSubClassOf(currentClass) {
		panic("java.lang.IllegalAccessError")
	}

	//mark 一般真的要执行的就是声明那个类的那个方法
	methodToBeInvoked := resolvedMethod

	//其实是父类方法
	if currentClass.IsSuper() &&//当前类的 ACC_SUPER标志被设置
		resolvedClass.IsSuperClassOf(currentClass) &&//超类中的函数
		resolvedMethod.Name() != "<init>"{//不是构造函数

		//需要一个额外的过程查找最终要调用的方法
		methodToBeInvoked = heap.LookupMethodInClass(
			currentClass.SuperClass(),
			methodRef.Name(),
			methodRef.Descriptor())
	}


	//如果查找过程失败，或者找到的方法是抽象的
	if methodToBeInvoked == nil || methodToBeInvoked.IsAbstract() {
		panic("java.lang.AbstractMethodError")
	}

	//执行
	base.InvokeMethod(frame, methodToBeInvoked)
}
