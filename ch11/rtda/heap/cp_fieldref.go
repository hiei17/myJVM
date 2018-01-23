package heap

import "jvmgo/ch11/classfile"
//运行时常量池里面的一项 字段引用 里面有:所属常量池 类名 自己的名字 描述符
type FieldRef struct {
	MemberRef//继承
	/*
		type MemberRef struct {
			SymRef
			name       string
			descriptor string
		}

		type SymRef struct {
			cp        *ConstantPool//本运行时常量池
			className string//引用的类名 字符串
			class     *Class//引用的类
		}
	*/

	field *Field//缓存解析后的字段指针 用到再去连接到字段本身
}

// classfile转class时用
func newFieldRef(cp *ConstantPool, refInfo *classfile.ConstantFieldrefInfo) *FieldRef {

	ref := &FieldRef{}
	ref.cp = cp//所属运行时常量池

	//父类方法
	// 复制:类名   字段名 描述符
	ref.copyMemberRefInfo(&refInfo.ConstantMemberrefInfo)
	return ref
}

//缓存了field
func (self *FieldRef) ResolvedField() *Field {
	if self.field == nil {
		self.resolveFieldRef()
	}
	return self.field
}

// jvms 5.4.3.2
func (self *FieldRef) resolveFieldRef() {

	thisClass := self.cp.class
	fieldInClass := self.ResolvedClass()
	field := lookupField(fieldInClass, self.name, self.descriptor)

	//2个错误
	if field == nil {
		panic("java.lang.NoSuchFieldError")
	}
	if !field.isAccessibleTo(thisClass) {
		panic("java.lang.IllegalAccessError")
	}

	self.field = field
}

func lookupField(c *Class, name, descriptor string) *Field {

	//本类找
	for _, field := range c.fields {
		if field.name == name && field.descriptor == descriptor {
			return field
		}
	}

	//接口找
	for _, iface := range c.interfaces {
		if field := lookupField(iface, name, descriptor); field != nil {
			return field
		}
	}

	//父类找
	if c.superClass != nil {
		return lookupField(c.superClass, name, descriptor)
	}

	return nil
}
