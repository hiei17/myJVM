package classfile

/*
CONSTANT_NameAndType_info {
    u1 tag;
    u2 name_index;
    u2 descriptor_index;
}
*/
/*Java语言支持方法重载（override），不同的方法可
以有相同的名字，只要参数列表不同即可。这就是为什么
CONSTANT_NameAndType_info结构要同时包含名称和描述符的原
因。*/
/*
Java是不能定义多个同名字段的，哪怕它们的类
型各不相同。这只是Java语法的限制而已，从class文件的层面来看，
是完全可以支持这点的。*/
//这个和类名一起 可以唯一确定一个方法或者字段
type ConstantNameAndTypeInfo struct {

	//2个都是常量池索引，指向CONSTANT_Utf8_info常量
	nameIndex       uint16//名称
	descriptorIndex uint16//描述符
}

func (self *ConstantNameAndTypeInfo) readInfo(reader *ClassReader) {
	self.nameIndex = reader.readUint16()
	self.descriptorIndex = reader.readUint16()
}

/*

//描述符
① byte、short、char、int、long、float double
 B、S、C、I、J、F和D
②引用类型的描述符是L＋类的完全限定名＋分号。
③数组类型的描述符是[＋数组元素类型描述符。
2）字段描述符就是字段类型的描述符。
3）方法描述符是（分号分隔的参数类型描述符）+返回值类型描述符，
void返回值由单个字母V表示。
*/
