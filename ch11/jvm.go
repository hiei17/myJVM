package main

import "fmt"
import "strings"
import "jvmgo/ch11/classpath"
import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/rtda/heap"

type JVM struct {
	cmd         *Cmd
	classLoader *heap.ClassLoader
	mainThread  *rtda.Thread
}

func newJVM(cmd *Cmd) *JVM {
	//加载路径 用户可以设置了
	//类加载器 它里面存着类路径 有加载.class的方法
	cp := classpath.Parse(cmd.XjreOption, cmd.cpOption)
	classLoader := heap.NewClassLoader(cp, cmd.verboseClassFlag)
	return &JVM{
		cmd:         cmd,
		classLoader: classLoader,
		mainThread:  rtda.NewThread(),
	}
}

func (self *JVM) start() {
	self.initVM()
	//执行主类的main（）方法
	self.execMain()
}

func (self *JVM) initVM() {
	//先加载sun.mis.VM类
	vmClass := self.classLoader.LoadClass("sun/misc/VM")
	base.InitClass(self.mainThread, vmClass)

	interpret(self.mainThread, self.cmd.verboseInstFlag)
}

//先加载主类，然后执行其main（）方法
func (self *JVM) execMain() {
	className := strings.Replace(self.cmd.class, ".", "/", -1)
	mainClass := self.classLoader.LoadClass(className)
	mainMethod := mainClass.GetMainMethod()
	if mainMethod == nil {
		fmt.Printf("Main method not found in class %s\n", self.cmd.class)
		return
	}
	//调用main（）方法之前，需要给它传递args参数，这是通过直接操作局部变量表实现的。
	argsArr := self.createArgsArray()
	frame := self.mainThread.NewFrame(mainMethod)
	frame.LocalVars().SetRef(0, argsArr)
	self.mainThread.PushFrame(frame)

	interpret(self.mainThread, self.cmd.verboseInstFlag)
}


//把Go的[]string变量转换成Java的字符串数组
//返回一个args []string 生成的java字符串数组对象
func (self *JVM) createArgsArray() *heap.Object {

	//得string类
	stringClass := self.classLoader.LoadClass("java/lang/String")

	//得string[]类
	arrStringClass := stringClass.ArrayClass()

	//数组容量
	arrCount := uint(len(self.cmd.args))
	//new 一个容量为arrCount 的string[]对象
	argsStringArrObject := arrStringClass.NewArray(arrCount)

	//string[]数组的数据
	objectData := argsStringArrObject.Refs()

	//把args一一放进去
	for i, arg := range self.cmd.args {

		//通过字符串池拿到这个字符串arg的字符串对象
		stringObject := heap.JString(self.classLoader, arg)
		//放入 字符串数组对象 里面
		objectData[i] = stringObject
	}
	return argsStringArrObject
}
