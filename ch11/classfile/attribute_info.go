package classfile

/*
常量是由Java虚拟机规范严格定义的，共有14种。
但属性是可以扩展的，不同的虚拟机实现可以定义自己的属性类型。
由于这个原因，Java虚拟机规范没有使用tag，
而是使用属性名来区别不同的属性。

每个属性机构
attribute_info {
    u2 attribute_name_index; //常量池索引，指向常量池中的CONSTANT_Utf8_info常量
    u4 attribute_length;
    u1 info[attribute_length];
}
*/


//接口
type AttributeInfo interface {

	//不同实现
	readInfo(reader *ClassReader)
}

//遍历 得到所有属性
func readAttributes(reader *ClassReader, cp ConstantPool) []AttributeInfo {

	attributesCount := reader.readUint16()

	attributes := make([]AttributeInfo, attributesCount)

	for i := range attributes {
		//一个属性
		attributes[i] = readAttribute(reader, cp)
	}
	return attributes
}

//解析一个属性 返回
func readAttribute(reader *ClassReader, cp ConstantPool) AttributeInfo {

	attrNameIndex := reader.readUint16()
	attrName := cp.getUtf8(attrNameIndex)

	attrLen := reader.readUint32()

	//返回不同实例 多态
	attrInfo := newAttributeInfo(attrName, attrLen, cp)
	//多态 不同实现读取
	attrInfo.readInfo(reader)
	return attrInfo
}

func newAttributeInfo(attrName string, attrLen uint32, cp ConstantPool) AttributeInfo {
	switch attrName {

		//method属性
		case "Code"://方法体
			return &CodeAttribute{cp: cp}
		//field属性
		case "ConstantValue"://常量表达式的值  存常量池坐标
			return &ConstantValueAttribute{}

		//method属性
		case "Exceptions"://变长属性，记录方法抛出的异常表
			return &ExceptionsAttribute{}

		//以下3是调试信息 不一定要 使用javac编译器编译Java程序时，默认会在class文件中生成 这些信息
		//  method属性的Code属性的属性
		case "LineNumberTable"://方法行号
			return &LineNumberTableAttribute{}
		//  method属性的Code属性的属性
		case "LocalVariableTable"://方法局部变量
			return &LocalVariableTableAttribute{}
		//class属性
		case "SourceFile"://源文件名 如 XXX.java
			return &SourceFileAttribute{cp: cp}

		//最简单的两种属性，仅起标记作用，不包含任何数据。
		//都可以用
		case "Synthetic"://为了支持嵌套类和嵌套接口
			return &SyntheticAttribute{}
		//都可以用
		case "Deprecated"://废弃标记
			return &DeprecatedAttribute{}

		//不支持
		default:
			return &UnparsedAttribute{attrName, attrLen, nil}
	}
}
