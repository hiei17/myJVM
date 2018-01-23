package main

import (
	"jvmgo/ch11/instructions/base"
	"fmt"

	"jvmgo/ch11/rtda"
	"jvmgo/ch11/instructions"

)





func interpret(thread *rtda.Thread, logInst bool) {
	defer catchErr(thread)
	loop(thread, logInst)
}




func loop(thread *rtda.Thread, logInst bool) {

	reader := &base.BytecodeReader{}
	for {
		frame := thread.CurrentFrame()
		pc := frame.NextPC()
		thread.SetPC(pc)
		// decode
		reader.Reset(frame.Method().Code(), pc)
		opcode := reader.ReadUint8()//字节码
		inst := instructions.NewInstruction(opcode)//指令对象
		inst.FetchOperands(reader)//拿到操作数
		frame.SetNextPC(reader.PC())

		//打印
		if (logInst) {
			logInstruction(frame, inst)
		}

		// execute
		inst.Execute(frame)

		//判断Java虚拟机栈中是否还有帧。如果没有则退出循环；否则继续。
		if thread.IsStackEmpty() {
			break
		}
	}
}



func catchErr(thread *rtda.Thread) {
	if r := recover(); r != nil {
		logFrames(thread)
		panic(r)
	}
}

//logFrames（）函数打印Java虚拟机栈信息
func logFrames(thread *rtda.Thread) {
	for !thread.IsStackEmpty() {
		frame := thread.PopFrame()
		method := frame.Method()
		className := method.Class().Name()
		fmt.Printf(">> pc:%4d %v.%v%v \n",
			frame.NextPC(), className, method.Name(), method.Descriptor())
	}
}

//logInstruction（）函数在方法执行过程中打印指令信息
func logInstruction(frame *rtda.Frame, inst base.Instruction) {
	method := frame.Method()
	className := method.Class().Name()
	methodName := method.Name()
	pc := frame.Thread().PC()
	fmt.Printf("%v.%v() #%2d %T %v\n", className, methodName, pc, inst, inst)
}