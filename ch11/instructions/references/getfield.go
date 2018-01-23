package references

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"

)

//获取对象的实例变量值，然后推入操作数栈
// Fetch field from object
type GET_FIELD struct{ base.Index16Instruction }

//栈顶实例弹出 实例.指定字段 的值 入栈
func (self *GET_FIELD) Execute(frame *rtda.Frame) {

	//todo 1 哪个实例变量
	field := frame.GetFieldByConstantIndex(self.Index)

	if field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	//todo 2 弹出栈顶实例
	stack := frame.OperandStack()
	//栈顶 是一个对象引用
	topObject := stack.PopRef()
	if topObject == nil {
		panic("java.lang.NullPointerException")
	}

	//todo 3.实例的这个变量的值 入栈
	fieldDescriptor := field.Descriptor()
	fieldSlotId := field.SlotId()
	topObjectAllFieldsSlots := topObject.Fields()

	switch fieldDescriptor[0] {
		case 'Z', 'B', 'C', 'S', 'I':
			stack.PushInt(topObjectAllFieldsSlots.GetInt(fieldSlotId))
		case 'F':
			stack.PushFloat(topObjectAllFieldsSlots.GetFloat(fieldSlotId))
		case 'J':
			stack.PushLong(topObjectAllFieldsSlots.GetLong(fieldSlotId))
		case 'D':
			stack.PushDouble(topObjectAllFieldsSlots.GetDouble(fieldSlotId))
		case 'L', '[':
			stack.PushRef(topObjectAllFieldsSlots.GetRef(fieldSlotId))
		default:
			// todo
	}
}
