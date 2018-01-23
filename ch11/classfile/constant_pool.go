package classfile

//常量池实际上也是一个表，

type ConstantPool [] ConstantInfo
//常量接口
type ConstantInfo interface {
	readInfo(reader *ClassReader)
}

//解析出常量池对象
func readConstantPool(reader *ClassReader) ConstantPool {

	cpCount := int(reader.readUint16())//常量个数

	cp := make([]ConstantInfo, cpCount)

	//表头给出的常量池大小比实际大1。假设表头给出的值是n，那么常量池的实际大小是n–1。，有效的常量池索引是1~n–1。0是无效索引，表示不指向任何常量。
	for i := 1; i < cpCount; i++ { //  注意索引从 1 开始
		cp[i] = readConstantInfo(reader, cp)
		switch cp[i].(type) {
		//CONSTANT_Long_info和 CONSTANT_Double_info各占两个位置
			case *ConstantLongInfo, *ConstantDoubleInfo:
				i++ //  占两个位置
		}
	}
	return cp
}


//解析出一个常量
func readConstantInfo(reader *ClassReader, cp ConstantPool) ConstantInfo {
	tag := reader.readUint8()//常量类型
	c := newConstantInfo(tag, cp)//多态实现ConstantInfo
	c.readInfo(reader)//读取常量信息，需要由具体的常量结构体实现
	return c
}

//工厂方法  得不同的常量对象
func newConstantInfo(tag uint8, cp ConstantPool) ConstantInfo {
	switch tag {

	//最小的数字 CONSTANT_Integer_info正好可以容纳一个Java的int型常量，
	//	但实际上比int更小的boolean、byte、short和char类型的常量也放在 CONSTANT_Integer_info中
	case CONSTANT_Integer: return &ConstantIntegerInfo{}
	case CONSTANT_Float: return &ConstantFloatInfo{}
	case CONSTANT_Long: return &ConstantLongInfo{}
	case CONSTANT_Double: return &ConstantDoubleInfo{}

	case CONSTANT_Utf8: return &ConstantUtf8Info{}
	case CONSTANT_String: return &ConstantStringInfo{cp: cp}
	case CONSTANT_Class: return &ConstantClassInfo{cp: cp}
	case CONSTANT_Fieldref:
		return &ConstantFieldrefInfo{ConstantMemberrefInfo{cp: cp}}
	case CONSTANT_Methodref:
		return &ConstantMethodrefInfo{ConstantMemberrefInfo{cp: cp}}
	case CONSTANT_InterfaceMethodref:
		return &ConstantInterfaceMethodrefInfo{ConstantMemberrefInfo{cp: cp}}
	case CONSTANT_NameAndType: return &ConstantNameAndTypeInfo{}
	case CONSTANT_MethodType: return &ConstantMethodTypeInfo{}
	case CONSTANT_MethodHandle: return &ConstantMethodHandleInfo{}
	case CONSTANT_InvokeDynamic: return &ConstantInvokeDynamicInfo{}
	default: panic("java.lang.ClassFormatError: constant pool tag!")
	}
}



/*****************************按索引从常量池get方法 不再是索引了 给直接的字符串***********************************/
//按索引查找常量
func (self ConstantPool) getConstantInfo(index uint16) ConstantInfo {
	if cpInfo := self[index]; cpInfo != nil {
		return cpInfo
	}
	panic("Invalid constant pool index!")
}
//常量池查找字段或方法的名字和描 述符
func (self ConstantPool) getNameAndType(index uint16) (string, string) {
	ntInfo := self.getConstantInfo(index).(*ConstantNameAndTypeInfo)
	name := self.getUtf8(ntInfo.nameIndex)
	_type := self.getUtf8(ntInfo.descriptorIndex)
	return name, _type
}
//从常量池查找类名
func (self ConstantPool) getClassName(index uint16) string {
	classInfo := self.getConstantInfo(index).(*ConstantClassInfo)
	return self.getUtf8(classInfo.nameIndex)
}
//从常量池查找UTF-8字符串
func (self ConstantPool) getUtf8(index uint16) string {
	utf8Info := self.getConstantInfo(index).(*ConstantUtf8Info)
	return utf8Info.str
}
/*****************************按索引从常量池get方法***************************************************/