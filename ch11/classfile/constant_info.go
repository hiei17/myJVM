package classfile
//常量
/*
数据结构
cp_info {
	u1 tag;
	u1 info[];
}
*/

//tag
//可以把常量池中的常量分为两类：
// 字面量（literal）和符号引用（symbolic reference）。
//字面量包括数字常量和字符串常量
const (
	//没做
	CONSTANT_MethodHandle = 15
	CONSTANT_MethodType = 16
	CONSTANT_InvokeDynamic = 18

	//3个都是 指向 CONSTANT_Class_info和CONSTANT_NameAndType_info
	CONSTANT_Fieldref = 9
	CONSTANT_Methodref = 10
	CONSTANT_InterfaceMethodref = 11

	//指向Utf8
	CONSTANT_String = 8
	CONSTANT_Class = 7//类名
	CONSTANT_NameAndType = 12//名字和描述符

	//最基本的常量
	CONSTANT_Utf8 = 1//字符串
	CONSTANT_Integer = 3//4字节 有符号 更小的int更小的boolean、byte、short和char类型 也是它
	CONSTANT_Float = 4//4字节有符号
	CONSTANT_Long = 5//8字节 有符号
	CONSTANT_Double = 6//8字节 有符号
)
