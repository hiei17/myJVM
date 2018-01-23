package rtda

import "jvmgo/ch11/rtda/heap"

//线程私有的内存
type Thread struct {
	pc int//计数器
	stack *Stack//指向栈顶
}

func NewThread() *Thread {
	return &Thread{
		//Java虚拟机栈可以是连续的空间，也可以不连续；
		// 可以是固定大小，也可以在运行时动态扩展 [1]  。
		// 如果Java虚拟机栈有大小限制， 且执行线程所需的栈空间超出了这个限制，会导致 StackOverflowError异常抛出。
		// 如果Java虚拟机栈可以动态扩展，但 是内存已经耗尽，会导致OutOfMemoryError异常抛出。
		stack: newStack(1024),//暂时设定为最多可以容纳多少1024栈帧Frame
		 }
}

func (self *Thread) NewFrame(method *heap.Method) *Frame {
	return newFrame(self, method)
}
func (self *Thread) PC() int { return self.pc } // getter
func (self *Thread) SetPC(pc int) { self.pc = pc } // setter
func (self *Thread) PushFrame(frame *Frame) {
	self.stack.push(frame)
}
func (self *Thread) PopFrame() *Frame {
	return self.stack.pop()
}
//返回当前帧 就是最顶上那个
func (self *Thread) CurrentFrame() *Frame {
	return self.stack.top()
}
func (self *Thread) TopFrame() *Frame {
	return self.stack.top()
}
func (self *Thread) IsStackEmpty() bool {
	return self.stack.isEmpty()
}

func (self *Thread) ClearStack() {
	self.stack.clear()
}


func (self *Thread) GetFrames() []*Frame {
	return self.stack.getFrames()
}