package rtda

// jvm stack
type Stack struct {
	maxSize uint//最大容量
	size    uint//当前放多少了
	//记最顶上就行了 每个都连下一个
	_top    *Frame // stack is implemented as linked list
}

func newStack(maxSize uint) *Stack {
	return &Stack{
		maxSize: maxSize,
	}
}

func (self *Stack) push(frame *Frame) {
	if self.size >= self.maxSize {
		//如果栈已经满了，按照Java虚拟机规范，应该抛出 StackOverflowError异常
		panic("java.lang.StackOverflowError")
	}

	if self._top != nil {
		//链?
		frame.lower = self._top
	}

	self._top = frame
	self.size++
}

func (self *Stack) pop() *Frame {
	if self._top == nil {
		panic("jvm stack is empty!")
	}

	top := self._top
	self._top = top.lower
	//拿出去的 是切断联系的
	top.lower = nil
	self.size--

	return top
}

func (self *Stack) top() *Frame {
	if self._top == nil {
		panic("jvm stack is empty!")
	}

	return self._top
}

func (self *Stack) isEmpty() bool {
	return self._top == nil
}

func (self *Stack) clear() {
	for !self.isEmpty() {
		self.pop()
	}
}


func (self *Stack) getFrames() []*Frame {
	frames := make([]*Frame, 0, self.size)
	for frame := self._top; frame != nil; frame = frame.lower {
		frames = append(frames, frame)
	}
	return frames
}
