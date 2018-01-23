package rtda

import (
	"math"
	"jvmgo/ch11/rtda/heap"
)

//局部变量表
type LocalVars []Slot


//局部变量表大小是编译器在编译时就确定的 在class文件的方法信息的code属性里面有写这项
func newLocalVars(maxLocals uint) LocalVars {
	if maxLocals > 0 {
		return make([]Slot, maxLocals)
	}
	return nil
}

//操作局部变量表和操作数栈的指令都是含类型信息的
//对boolean、byte、short和char类型 , 都可以转换成int值类来处理
//最基础的
func (self LocalVars) SetInt(index uint, val int32) {
	self[index].num = val
}
func (self LocalVars) GetInt(index uint) int32 {
	return self[index].num
}
//float变量可以先转成int类型，然后按int变量来处理。
func (self LocalVars) SetFloat(index uint, val float32) {
	bits := math.Float32bits(val)
	self[index].num = int32(bits)
}
func (self LocalVars) GetFloat(index uint) float32 {
	bits := uint32(self[index].num)
	return math.Float32frombits(bits)
}
//long变量则需要拆成两个int变量。
func (self LocalVars) SetLong(index uint, val int64) {
	self[index].num = int32(val)
	self[index+1].num = int32(val >> 32)
}
func (self LocalVars) GetLong(index uint) int64 {
	low := uint32(self[index].num)
	high := uint32(self[index+1].num)
	return int64(high)<<32 | int64(low)
}
//double变量可以先转成long类型，然后按照long变量来处理。
func (self LocalVars) SetDouble(index uint, val float64) {
	bits := math.Float64bits(val)
	self.SetLong(index, int64(bits))
}
func (self LocalVars) GetDouble(index uint) float64 {
	bits := uint64(self.GetLong(index))
	return math.Float64frombits(bits)
}
//最后是引用值，也比较简单，直接存取即可。
func (self LocalVars) SetRef(index uint, ref *heap.Object) {
	self[index].ref = ref
}
func (self LocalVars) GetRef(index uint) *heap.Object {
	return self[index].ref
}
func (self LocalVars) GetThis() *heap.Object {
	return self.GetRef(0)
}
func (self LocalVars) SetSlot(index uint, slot Slot) {
	self[index] = slot
}