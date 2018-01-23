package heap

import "jvmgo/ch11/classfile"

//访问标志 名字 描述符 所属类 值在常量池的索引
type Field struct {

	ClassMember //和方法公共的
	/*type ClassMember struct {
				  //3个都直接从classfile.MemberInfo里面拿
		accessFlags uint16
		name        string
		descriptor  string

				  //Class结构体指针，这样可以通过字段或方法访问到它所属的类。
				  //构造函数传进来赋值
		class       *Class//所属类
	}*/
	slotId          uint//字段的值 存储位置的编号
	constValueIndex uint//值的常量池编号
}


//从classfile.MemberInfo产生字段
//new 出来后每个字段有:
//	访问标志 名字 描述  值在常量池的索引
func newFields(class *Class, cfFields []*classfile.MemberInfo) []*Field {

	fields := make([]*Field, len(cfFields))

	//遍历classfile中所有 字段
	for i, cfField := range cfFields {

		fields[i] = &Field{}
		fields[i].class = class//所属类

		//父类方法
		// 复制: 访问标志 名字 描述
		fields[i].copyMemberInfo(cfField)

		//字段值的 //常量池索引，但具体指向哪种常量因字段类型而异。
		fields[i].copyValueIndexAttributes(cfField)
	}
	return fields
}


//找到classfile.MemberInfo 的常量属性(常量池坐标) 返回
func (self *Field) copyValueIndexAttributes(cfField *classfile.MemberInfo) {
	valAttr := cfField.ConstantValueAttribute();
	if  valAttr != nil {
		self.constValueIndex = uint(valAttr.ConstantValueIndex())
	}
}



func (self *Field) IsVolatile() bool {
	return 0 != self.accessFlags&ACC_VOLATILE
}
func (self *Field) IsTransient() bool {
	return 0 != self.accessFlags&ACC_TRANSIENT
}
func (self *Field) IsEnum() bool {
	return 0 != self.accessFlags&ACC_ENUM
}

func (self *Field) ConstValueIndex() uint {
	return self.constValueIndex
}
func (self *Field) SlotId() uint {
	return self.slotId
}


/*// reflection
func (self *Field) Type() *Class {
	className := toClassName(self.descriptor)
	return self.class.loader.LoadClass(className)
}*/
func (self *Field) isLongOrDouble() bool {
	return self.descriptor == "J" || self.descriptor == "D"
}