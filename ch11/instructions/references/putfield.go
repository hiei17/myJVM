package references

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/rtda/heap"

// Set field in object //指令有一个操作数 是常量池坐标 可得字段
type PUT_FIELD struct{ base.Index16Instruction }


//字节码指定字段  栈里面有 要赋的值 和字段所在实例
func (self *PUT_FIELD) Execute(frame *rtda.Frame) {

	currentMethod := frame.Method()//创造本帧的方法 有存着引用
	currentClass := currentMethod.Class()//方法是类成员 里面存着类的引用
	cp := currentClass.ConstantPool()//类的运行时常量池
	fieldRef := cp.GetConstant(self.Index).(*heap.FieldRef)//字节码的中得到的操作数 就是坐标
	//字段对象引用里面 有存着字段所在类 描述 名字   得字段对象
	field := fieldRef.ResolvedField()

	//不是类变量才行
	if field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	//要是final 赋值必须在构造方法里面
	if field.IsFinal() {
		if currentClass != field.Class() || currentMethod.Name() != "<init>" {
			panic("java.lang.IllegalAccessError")
		}
	}

	descriptor := field.Descriptor()
	//得这个field在类的slots里面的编号
	slotId := field.SlotId()

	stack := frame.OperandStack()

	switch descriptor[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		val := stack.PopInt()//数
		ref := stack.PopRef()//实例
		//如果引用是null，需要抛出著名的空指针异常
		if ref == nil {
			panic("java.lang.NullPointerException")
		}
		ref.Fields().SetInt(slotId, val)

	case 'F':
		val := stack.PopFloat()
		ref := stack.PopRef()
		if ref == nil {
			panic("java.lang.NullPointerException")
		}
		ref.Fields().SetFloat(slotId, val)
	case 'J':
		val := stack.PopLong()
		ref := stack.PopRef()
		if ref == nil {
			panic("java.lang.NullPointerException")
		}
		ref.Fields().SetLong(slotId, val)
	case 'D':
		val := stack.PopDouble()
		ref := stack.PopRef()
		if ref == nil {
			panic("java.lang.NullPointerException")
		}
		ref.Fields().SetDouble(slotId, val)
	case 'L', '[':
		val := stack.PopRef()
		ref := stack.PopRef()
		if ref == nil {
			panic("java.lang.NullPointerException")
		}
		ref.Fields().SetRef(slotId, val)
	default:
		// todo
	}

}
