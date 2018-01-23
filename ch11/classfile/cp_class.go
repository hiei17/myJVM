package classfile

/*
CONSTANT_Class_info {
    u1 tag;
    u2 name_index;
}
*/
//表示类或者接口的符号引用
type ConstantClassInfo struct {
	cp        ConstantPool
	nameIndex uint16//常量池索引，指 向CONSTANT_Utf8_info常量 那里面是全限定名的.变/ 类名
}

func (self *ConstantClassInfo) readInfo(reader *ClassReader) {
	self.nameIndex = reader.readUint16()
}

//get到的就不是索引了  直接是字符串
func (self *ConstantClassInfo) Name() string {
	return self.cp.getUtf8(self.nameIndex)
}
