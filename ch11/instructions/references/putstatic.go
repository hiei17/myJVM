package references

import (
	"jvmgo/ch11/rtda"
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda/heap"

)
/*

putstatic指令:类变量赋值
 被赋值的类字段 由指令操作数(常量池field引用项)指定
值在栈顶
*/

//出栈给类变量赋值
// Set static field in class
type PUT_STATIC struct{ base.Index16Instruction }




//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (self *PUT_STATIC) Execute(frame *rtda.Frame) {
	currentMethod := frame.Method()//当前方法 //创造本帧的方法
	currentClass := currentMethod.Class()//当前类
	constantPool := currentClass.ConstantPool()//当前常量池
	fieldRef := constantPool.GetConstant(self.Index).(*heap.FieldRef)

	//字段对象引用里面 有存着字段所在类 描述 名字   得字段对象
	field := fieldRef.ResolvedField()//解析字段符号引
	class := field.Class()

	//mark 没初始化
	if !class.InitStarted() {
		frame.RevertNextPC()
		base.InitClass(frame.Thread(), class)
		return
	}

	//字段是实例字段而非静态字段
	if !field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	//如果是final字段，则实际操作的 是静态常量，
	if field.IsFinal() {
		//只能在类初始化方法中给它赋值
		if currentClass != class || currentMethod.Name() != "<clinit>" {//类初始化方法由编译器生成，名字是 <clinit>
			panic("java.lang.IllegalAccessError")
		}
	}
	descriptor := field.Descriptor()//写着数据类型
	slotId := field.SlotId()//这个类变量在 本类所有类变量的数组的索引
	slots := class.StaticVars()//本类所有类变量的数组
	stack := frame.OperandStack()//操作数栈

	//字段类型从操作数栈中弹出相应的值，然后赋给静态变量。  //不同类型的数据 不同set pop

	switch descriptor[0] {
		case 'Z', 'B', 'C', 'S', 'I': slots.SetInt(slotId, stack.PopInt())
		case 'F': slots.SetFloat(slotId, stack.PopFloat())
		case 'J': slots.SetLong(slotId, stack.PopLong())
		case 'D': slots.SetDouble(slotId, stack.PopDouble())
		case 'L', '[':
			ref := stack.PopRef()
			/*if(ref==nil) {
				panic("putstaticNil:"+class.Name()+"的"+fieldRef.Name())
			}*/
			slots.SetRef(slotId, ref)
	}
}
