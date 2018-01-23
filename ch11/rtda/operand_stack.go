package rtda
//操作数栈 属于栈帧 一个方法执行就对应一个
import (
	"math"
	"jvmgo/ch11/rtda/heap"
)
type OperandStack struct {
	size uint
	slots []Slot
}

//操作数栈 大小是
// 编译器在编译时就确定的 在class文件的方法信息的code属性里面有写这项
func newOperandStack(maxStack uint) *OperandStack {
	if maxStack > 0 {
		return &OperandStack{
			slots: make([]Slot, maxStack),
		}
	}
	return nil
}

//栈顶放一个int变量，然后把size加1
func (self *OperandStack) PushInt(val int32) {
	self.slots[self.size].num = val
	self.size++
}
//先把size减1，然后返回变量值
func (self *OperandStack) PopInt() int32 {
	self.size--
	return self.slots[self.size].num
}
//float变量还是先 转成int类型，然后按int变量处理。
func (self *OperandStack) PushFloat(val float32) {
	bits := math.Float32bits(val)
	self.slots[self.size].num = int32(bits)
	self.size++
}
func (self *OperandStack) PopFloat() float32 {
	self.size--
	bits := uint32(self.slots[self.size].num)
	return math.Float32frombits(bits)
}
//把long变量推入栈顶时，要拆成两个int变量。
func (self *OperandStack) PushLong(val int64) {
	self.slots[self.size].num = int32(val)
	self.slots[self.size+1].num = int32(val >> 32)
	self.size += 2
}
//弹出时，先弹出 两个int变量，然后组装成一个long变量。
func (self *OperandStack) PopLong() int64 {
	self.size -= 2
	low := uint32(self.slots[self.size].num)
	high := uint32(self.slots[self.size+1].num)
	return int64(high)<<32 | int64(low)
}
//double变量先转成long类型，然后按long变量处理。
func (self *OperandStack) PushDouble(val float64) {
	bits := math.Float64bits(val)
	self.PushLong(int64(bits))
}
func (self *OperandStack) PopDouble() float64 {
	bits := uint64(self.PopLong())
	return math.Float64frombits(bits)
}

func (self *OperandStack) PushSlot(slot Slot) {
	self.slots[self.size] = slot
	self.size++
}
func (self *OperandStack) PopSlot() Slot {
	self.size--
	return self.slots[self.size]
}
func (self *OperandStack) PushBoolean(val bool) {
	if val {
		self.PushInt(1)
	} else {
		self.PushInt(0)
	}
}


func (self *OperandStack) PushRef(ref *heap.Object) {
	self.slots[self.size].ref = ref
	self.size++
}
func (self *OperandStack) PopRef() *heap.Object {
	self.size--
	ref := self.slots[self.size].ref

	//todo  是为了帮助Go的垃圾收集器回收Object结构体实例。
	self.slots[self.size].ref = nil
	return ref
}

func (self *OperandStack) GetRefFromTop(n uint) *heap.Object {
	return self.slots[self.size-1-n].ref
}


func (self *OperandStack) Clear() {
	self.size = 0
	for i := range self.slots {
		self.slots[i].ref = nil
	}
}
func NewOperandStack(maxStack uint) *OperandStack {
	return newOperandStack(maxStack)
}
