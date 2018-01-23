package references

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"

//超数组指定field 找到它的class  把这个static field值取出 入栈

//类变量入栈 
// Get static field from class
type GET_STATIC struct{ base.Index16Instruction }



func (self *GET_STATIC) Execute(frame *rtda.Frame)  {
	
	//mark 1.哪个类变量
	field:=frame.GetFieldByConstantIndex(self.Index)

	//不是static 报异常
	if !field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}
	class := field.Class()

	//mark 没初始化
	if !class.InitStarted() {
		frame.RevertNextPC()
		base.InitClass(frame.Thread(), class)
		return
	}

	descriptor := field.Descriptor()
	slotId := field.SlotId()//坐标
	slots := class.StaticVars()//这类静态变量
	
	
	stack := frame.OperandStack()//操作数栈

	//todo 2.得到类变量值 压入操作数栈
	//不同类型的入栈

	switch descriptor[0] {
		case 'Z', 'B', 'C', 'S', 'I':
			stack.PushInt(slots.GetInt(slotId))
		case 'F':
			stack.PushFloat(slots.GetFloat(slotId))
		case 'J':
			stack.PushLong(slots.GetLong(slotId))
		case 'D':
			stack.PushDouble(slots.GetDouble(slotId))
		case 'L', '[':
			ref := slots.GetRef(slotId)
			/*if(ref==nil){
				panic(class.Name()+"的"+field.Name()+"没有")
			}*/
			stack.PushRef(ref)

		default:
			// todo
	}
}
