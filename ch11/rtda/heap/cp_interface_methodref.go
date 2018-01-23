package heap

import "jvmgo/ch11/classfile"

type InterfaceMethodRef struct {
	MemberRef
	/*
		type MemberRef struct {
			SymRef
			name       string
			descriptor string
		}

		type SymRef struct {
			cp        *ConstantPool//本运行时常量池
			className string//引用的类名 //常量池索引，指 向CONSTANT_Utf8_info常量
			class     *Class//引用的类
		}
	*/
	method *Method
}

//classfile生成运行时常量池的时候调用
func newInterfaceMethodRef(cp *ConstantPool, refInfo *classfile.ConstantInterfaceMethodrefInfo) *InterfaceMethodRef {

	ref := &InterfaceMethodRef{}
	ref.cp = cp

	//复制  所属类名 名字 描述符
	ref.copyMemberRefInfo(&refInfo.ConstantMemberrefInfo)
	return ref
}


//接口方法符号引用的解析
func (self *InterfaceMethodRef) ResolvedInterfaceMethod() *Method {
	if self.method == nil {
		//class名 方法名 描述符 这3坐标对上的方法就是要找的方法
		self.resolveInterfaceMethodRef()
	}
	return self.method
}

// jvms8 5.4.3.4
func (self *InterfaceMethodRef) resolveInterfaceMethodRef() {
	myClass := self.cp.class
	methodInClass := self.ResolvedClass()

	//所在类是接口
	if !methodInClass.IsInterface() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	method := lookupInterfaceMethod(methodInClass, self.name, self.descriptor)

	if method == nil {
		panic("java.lang.NoSuchMethodError")
	}
	if !method.isAccessibleTo(myClass) {
		panic("java.lang.IllegalAccessError")
	}

	self.method = method
}

// 在本接口中找 找不到再去本接口的接口找
func lookupInterfaceMethod(iface *Class, name, descriptor string) *Method {
	for _, method := range iface.methods {
		if method.name == name && method.descriptor == descriptor {
			return method
		}
	}

	return lookupMethodInInterfaces(iface.interfaces, name, descriptor)
}
