package references

import "jvmgo/ch11/instructions/base"
import "jvmgo/ch11/rtda"
import "jvmgo/ch11/rtda/heap"

//int[] a1 = new int[10]; 创建基本类型数组 入栈

const (
	//Array Type  atype
	AT_BOOLEAN = 4
	AT_CHAR    = 5
	AT_FLOAT   = 6
	AT_DOUBLE  = 7
	AT_BYTE    = 8
	AT_SHORT   = 9
	AT_INT     = 10
	AT_LONG    = 11
)

// Create new array 操作数是数组类型
type NEW_ARRAY struct {
	atype uint8
}

func (self *NEW_ARRAY) FetchOperands(reader *base.BytecodeReader) {
	self.atype = reader.ReadUint8()
}

func (self *NEW_ARRAY) Execute(frame *rtda.Frame) {

	stack := frame.OperandStack()
	count := stack.PopInt()//数组容量 操作数栈弹出
	if count < 0 {
		panic("java.lang.NegativeArraySizeException")
	}

	classLoader := frame.Method().Class().Loader()
	//指定数据类型 让类加载器生成数组类
	arrClass := getPrimitiveArrayClass(classLoader, self.atype)
	//给数组大小 和 数组类 生成数组
	arr := arrClass.NewArray(uint(count))
	stack.PushRef(arr)
}

func getPrimitiveArrayClass(loader *heap.ClassLoader, atype uint8) *heap.Class {
	switch atype {
		case AT_BOOLEAN:
			return loader.LoadClass("[Z")
		case AT_BYTE:
			return loader.LoadClass("[B")
		case AT_CHAR:
			return loader.LoadClass("[C")
		case AT_SHORT:
			return loader.LoadClass("[S")
		case AT_INT:
			return loader.LoadClass("[I")
		case AT_LONG:
			return loader.LoadClass("[J")
		case AT_FLOAT:
			return loader.LoadClass("[F")
		case AT_DOUBLE:
			return loader.LoadClass("[D")
		default:
			panic("Invalid atype!")
	}
}
