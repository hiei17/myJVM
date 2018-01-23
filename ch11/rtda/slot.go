package rtda
//局部变量表的slot
import "jvmgo/ch11/rtda/heap"
//局部变量表是按索引访问的，根据Java虚拟机规范，这个数组的每个元素至少可以容纳一个int或引用值，
//两个连续的元素可以容纳一个long或double值。
type Slot struct {
	num int32//整形
	ref *heap.Object//引用
}
