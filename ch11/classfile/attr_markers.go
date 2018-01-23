package classfile

/*
Deprecated_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
}
*/
type DeprecatedAttribute struct {
	MarkerAttribute
}

/*
Synthetic_attribute {
    u2 attribute_name_index;
    u4 attribute_length;//必须是0
}
*///仅起标记作用
type SyntheticAttribute struct {
	MarkerAttribute
}

//空的
type MarkerAttribute struct{}

//空的
func (self *MarkerAttribute) readInfo(reader *ClassReader) {
	// read nothing
}
