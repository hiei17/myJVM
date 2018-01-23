package heap

import "jvmgo/ch11/classfile"

//运行时常量池的一个项目 表对一个方法放引用
//里面有:类名 方法名 描述符
type MethodRef struct {
	MemberRef
	/*
	cp        *ConstantPool//本运行时常量池
	className string//引用的类名 字符串
	class     *Class//引用的类
	name       string //方法名 字符串
	descriptor string//描述符字符串
	*/
	method *Method
}

//classfile生成运行时常量池时调用
func newMethodRef(cp *ConstantPool, refInfo *classfile.ConstantMethodrefInfo) *MethodRef {
	ref := &MethodRef{}
	ref.cp = cp
	//父类方法 从classfile里面得到本成员引用的 类名 名字 描述
	ref.copyMemberRefInfo(&refInfo.ConstantMemberrefInfo)
	return ref
}

//方法符号引用的解析
func (self *MethodRef) ResolvedMethod() *Method {
	if self.method == nil {
		//还没有解析过符号引用
		self.resolveMethodRef()
	}
	return self.method
}

//解析出引用的method
// jvms8 5.4.3.3
func (self *MethodRef) resolveMethodRef() {
	//本类
	myClass := self.cp.class

	//引用的类
	methodInClass := self.ResolvedClass()

	//调用的方法在接口里面  不可能 IncompatibleClassChangeError异常，
	if methodInClass.IsInterface() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	method := lookupMethod(methodInClass, self.name, self.descriptor)

	if method == nil {
		panic("java.lang.NoSuchMethodError")
	}

	//检查访问权限
	if !method.isAccessibleTo(myClass) {
		panic("java.lang.IllegalAccessError")
	}

	self.method = method
}


//继承层次中找，如果找不到，就去它的接口中找。
func lookupMethod(class *Class, name, descriptor string) *Method {

	//继承
	method := LookupMethodInClass(class, name, descriptor)

	if method == nil {
		//接口
		method = lookupMethodInInterfaces(class.interfaces, name, descriptor)
	}
	return method
}


