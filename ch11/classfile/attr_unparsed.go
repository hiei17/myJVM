package classfile

/*
//不支持的属性
attribute_info {
    u2 attribute_name_index;
    u4 attribute_length;
    u1 info[attribute_length];//自己返回没解析的内容
}
*/
type UnparsedAttribute struct {
	name   string
	length uint32
	info   []byte
}

//接口要求的方法
func (self *UnparsedAttribute) readInfo(reader *ClassReader) {
	self.info = reader.readBytes(self.length)
}

//get方法
func (self *UnparsedAttribute) Info() []byte {
	return self.info
}
