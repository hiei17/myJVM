package classfile

/*
SourceFile_attribute {
    u2 attribute_name_index;
    u4 attribute_length;//长度必须是2
    u2 sourcefile_index;
}
*/
type SourceFileAttribute struct {
	cp              ConstantPool
	sourceFileIndex uint16 //是常量池索引， 指向CONSTANT_Utf8_info常量
}

func (self *SourceFileAttribute) readInfo(reader *ClassReader) {
	self.sourceFileIndex = reader.readUint16()
}

//get方法 返回源文件名 如"ClassFileTest.java"
func (self *SourceFileAttribute) FileName() string {

	return self.cp.getUtf8(self.sourceFileIndex)
}
