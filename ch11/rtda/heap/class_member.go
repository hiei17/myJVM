package heap

import "jvmgo/ch11/classfile"


//方法和字段 都是成员
type ClassMember struct {
			  //3个都直接从classfile.MemberInfo里面拿
	accessFlags uint16
	name        string
	descriptor  string

	//Class结构体指针，这样可以通过字段或方法访问到它所属的类。
			  //构造函数传进来赋值
	class       *Class//所属类
}

//从classfile.MemberInfo 拿来 访问标志 名字 描述
func (self *ClassMember) copyMemberInfo(memberInfo *classfile.MemberInfo) {
	self.accessFlags = memberInfo.AccessFlags()
	//不是常量池索引了 拿出来了
	self.name = memberInfo.Name()
	self.descriptor = memberInfo.Descriptor()
}


func (self *ClassMember) IsPublic() bool {
	return 0 != self.accessFlags&ACC_PUBLIC
}
func (self *ClassMember) IsPrivate() bool {
	return 0 != self.accessFlags&ACC_PRIVATE
}
func (self *ClassMember) IsProtected() bool {
	return 0 != self.accessFlags&ACC_PROTECTED
}
func (self *ClassMember) IsStatic() bool {
	return 0 != self.accessFlags&ACC_STATIC
}
func (self *ClassMember) IsFinal() bool {
	return 0 != self.accessFlags&ACC_FINAL
}
func (self *ClassMember) IsSynthetic() bool {
	return 0 != self.accessFlags&ACC_SYNTHETIC
}

// getters
func (self *ClassMember) Name() string {
	return self.name
}
func (self *ClassMember) Descriptor() string {
	return self.descriptor
}
func (self *ClassMember) Class() *Class {
	return self.class
}

// jvms 5.4.4
func (refMember *ClassMember) isAccessibleTo(thisClass *Class) bool {
	//Public怎么样都能拿到
	if refMember.IsPublic() {
		return true
	}

	memberInClass := refMember.class

	//Protected 子类和同包能拿到
	if refMember.IsProtected() {
		return thisClass == memberInClass || //同类
			thisClass.IsSubClassOf(memberInClass) ||//继承
			memberInClass.GetPackageName() == thisClass.GetPackageName()//同包
	}
	//default 同包能拿到
	if !refMember.IsPrivate() {
		return memberInClass.GetPackageName() == thisClass.GetPackageName()
	}

	//private
	return thisClass == memberInClass
}

