package comparisons

import (
	"jvmgo/ch11/instructions/base"
	"jvmgo/ch11/rtda"
)


/*

比较指令可以分为两类：
一类将比较结果推入操作数栈顶，
一类根据比较结果跳转。
比较指令是编译器实现if-else、for、while等语句的基石，共有19条
*/

// Compare long
type LCMP struct{ base.NoOperandsInstruction }
/*栈顶的两个long变量弹出，进行比较，
比较结果（int型0、1或-1）推入栈顶*/
func (self *LCMP) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopLong()
	v1 := stack.PopLong()
	if v1 > v2 {
		stack.PushInt(1)
	} else if v1 == v2 {
		stack.PushInt(0)
	} else {
		stack.PushInt(-1)
	}
}
