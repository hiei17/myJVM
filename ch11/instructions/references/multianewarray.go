package references

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"
import (
	"jvmgo/ch11/rtda/heap"
)

//a=new int[3][4][5]; 入栈

// Create new multidimensional array
type MULTI_ANEW_ARRAY struct {
	index      uint16//常量池索引
	dimensions uint8//维度数
}

func (self *MULTI_ANEW_ARRAY) FetchOperands(reader *base.BytecodeReader) {
	self.index = reader.ReadUint16()
	self.dimensions = reader.ReadUint8()
}
func (self *MULTI_ANEW_ARRAY) Execute(frame *rtda.Frame) {
	//多维数组本身类型(和其他数组不同
	cp := frame.Method().Class().ConstantPool()
	classRef := cp.GetConstant(uint(self.index)).(*heap.ClassRef)
	multiClass := classRef.ResolvedClass()//int[][][]类

	stack := frame.OperandStack()
	//弹出维度数个整数，分别代表每一个维度的数组长度。
	counts := popAndCheckCounts(stack, int(self.dimensions))//{3,4,5}
	//数组对象
	arrObject := newMultiDimensionalArray(counts, multiClass)
	stack.PushRef(arrObject)
}

func popAndCheckCounts(stack *rtda.OperandStack, dimensions int) []int32 {
	counts := make([]int32, dimensions)
	for i := dimensions - 1; i >= 0; i-- {
		counts[i] = stack.PopInt()
		if counts[i] < 0 {
			panic("java.lang.NegativeArraySizeException")
		}
	}

	return counts
}

func newMultiDimensionalArray(counts []int32, arrClass *heap.Class) *heap.Object {
	count := uint(counts[0])//3
	arr := arrClass.NewArray(count)//a[3]实例

	if len(counts) > 1 {
		refs := arr.Refs()
		for i := range refs {
			//数组组成元素的类
			componentClass := arrClass.ComponentClass()//int[][]
			//a[0]=int[4][]
			refs[i] = newMultiDimensionalArray(counts[1:], componentClass)
		}
	}

	return arr
}
