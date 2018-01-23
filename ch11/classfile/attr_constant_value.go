package classfile

/*
ConstantValue_attribute {
    u2 attribute_name_index;//ConstantValue
    u4 attribute_length;//必须2
    u2 constantvalue_index;//常量表达式的值
}
*/
type ConstantValueAttribute struct {
	constantValueIndex uint16
	//常量池索引，但具体指向哪种常量因字段类型而异。
}

func (self *ConstantValueAttribute) readInfo(reader *ClassReader) {
	self.constantValueIndex = reader.readUint16()
}

func (self *ConstantValueAttribute) ConstantValueIndex() uint16 {
	return self.constantValueIndex
}
