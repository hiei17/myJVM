package heap

import "jvmgo/ch11/classfile"
//在运行时常量池里面记载的类引用

type ClassRef struct {
	SymRef
	/*
		type SymRef struct {
			cp        *ConstantPool//本运行时常量池
			className string//引用的类名
			class     *Class//引用的类
		}
	*/
}

//传参:所在的运行时常量池  常量
//new 出来 只有运行时常量池 和 类名
func newClassRef(cp *ConstantPool, classInfo *classfile.ConstantClassInfo) *ClassRef {
	//new 一个空的类引用常量
	ref := &ClassRef{}
	ref.cp = cp//所在的运行时常量池
	ref.className = classInfo.Name()//类的全名 直接是字符串
	return ref
}
