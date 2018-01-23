package heap

import "jvmgo/ch11/classfile"

//方法和字段的父类
//共有部分是:所属常量池 类名 自己的名字 描述符
type MemberRef struct {
	SymRef//继承
	/*type SymRef struct {
		cp        *ConstantPool//本运行时常量池
		className string//引用的类名 字符串
		class     *Class//引用的类 本来是空的 用到再加载
	}*/
	name       string
	descriptor string
}

//从classfile里面复制本成员的 类名 名字 描述
func (self *MemberRef) copyMemberRefInfo(refInfo *classfile.ConstantMemberrefInfo) {
	//类名字符串
	self.className = refInfo.ClassName()
	//成员名 描述符 的字符串
	self.name, self.descriptor = refInfo.NameAndDescriptor()
}


////////////////////////////////////////////////////////get方法///////////////////////////////////////////////////////////
func (self *MemberRef) Name() string {
	return self.name
}
func (self *MemberRef) Descriptor() string {
	return self.descriptor
}
