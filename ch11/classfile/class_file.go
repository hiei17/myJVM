package classfile

import "fmt"


//.class文件的直接1:1描述形成的数据结构(完全未经加工  和Classpy里面看到的一样
// 除了常量池项tag不用写了  属性名不用写了 用结构体类型了 不用表用数组了也不用表头了
// 作为生产class的中间层

/*
ClassFile {
    u4             magic;
    u2             minor_version;
    u2             major_version;
    u2             constant_pool_count;
    cp_info        constant_pool[constant_pool_count-1];
    u2             access_flags;
    u2             this_class;
    u2             super_class;
    u2             interfaces_count;
    u2             interfaces[interfaces_count];
    u2             fields_count;
    field_info     fields[fields_count];
    u2             methods_count;
    method_info    methods[methods_count];
    u2             attributes_count;
    attribute_info attributes[attributes_count];
}
*/
type ClassFile struct {
	//magic      uint32 魔数检测下是不是这个0xCAFEBABE就行 不用存
	minorVersion uint16//00 00
	majorVersion uint16//00 34
	constantPool ConstantPool//内部是数组
	accessFlags  uint16
	thisClass    uint16
	superClass   uint16
	interfaces   []uint16
	fields       []*MemberInfo
	methods      []*MemberInfo
	attributes   []AttributeInfo
}

//[]byte → ClassFile结构体
//是个函数 函数名(入参)(返回值)
func Parse(classData []byte) (cf *ClassFile, err error) {

	//Go语言没有异常处理机制，只有一个panic-recover机制
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	cr := &ClassReader{classData}//byte数组解析工具
	cf = &ClassFile{}//new class文件对象
	cf.read(cr)//用ClassReader解析出一个ClassFile对象cf
	return
}

//依次调用其他方法解析class文件
// 是一个方法 改变在self 也就是一个ClassFile对象
//func(对象)方法名(入参)
func (self *ClassFile) read(reader *ClassReader) {

	self.readAndCheckMagic(reader)//CAFEBABE
	self.readAndCheckVersion(reader)//00 00 00 34

	//常量池
	self.constantPool = readConstantPool(reader)

	//访问标志
	//16位的“bitmask”，指 出class文件定义的是类还是接口，访问级别是public还是private，
	// 这里只对class文件进行初步解析， 只是 读取类访问标志以备后用。
	self.accessFlags = reader.readUint16()//00 21

	//常量池坐标 只是把全限定名 . 换 /
	//本类
	self.thisClass = reader.readUint16()//00 05//常量池坐标
	//父类
	self.superClass = reader.readUint16()//00 06//常量池坐标
	//接口索引表，表中存放的也是常量池索引，给出该类实现的所有接口的名字。
	self.interfaces = reader.readUint16s()//00 00 没有接口


	//字段和方法的结相同，差别仅在于属性表  这里用一种数据结构存
	self.fields = readMembers(reader, self.constantPool)//字段数组
	self.methods = readMembers(reader, self.constantPool)//方法数组

	//属性
	/*
		attribute_info {
		    u2 attribute_name_index; //常量池索引，指向常量池中的CONSTANT_Utf8_info常量
		    u4 attribute_length;
		    u1 info[attribute_length];
		}
	*/
	self.attributes = readAttributes(reader, self.constantPool)

}

//检测魔数
func (self *ClassFile) readAndCheckMagic(reader *ClassReader) {
	magic := reader.readUint32()//0xCAFEBABE  4个字节 每个字节可用2个16进制数表示

	//如果加载的class文件不符合要求的格式
	// 就抛出java.lang.ClassFormatError异常
	if magic != 0xCAFEBABE {
		//暂时先调用 panic（）方法终止程序执行
		panic("java.lang.ClassFormatError: magic!")
	}
}

//检测版本号
func (self *ClassFile) readAndCheckVersion(reader *ClassReader) {
	self.minorVersion = reader.readUint16()//2个字节 00 00
	self.majorVersion = reader.readUint16()//2个字节 00 34
	switch self.majorVersion {
	case 45:
		return
	case 46, 47, 48, 49, 50, 51, 52://jdk8是52  16进制是34
		//次版本号后来都没用了 都是00了
		if self.minorVersion == 0 {
			return
		}
	}

	//如果遇到其他版本号， 暂时先调用panic（）方法终止程序执行
	panic("java.lang.UnsupportedClassVersionError!")
}

//go里面大写字母开头的方法字段结构体 就是包外也可访问的 小写只有用包能访问
//这6个简单的get方法 公开给任何包调用  简单返回相应字段
func (self *ClassFile) MinorVersion() uint16 {
	return self.minorVersion
}
func (self *ClassFile) MajorVersion() uint16 {
	return self.majorVersion
}
func (self *ClassFile) ConstantPool() ConstantPool {
	return self.constantPool
}
func (self *ClassFile) AccessFlags() uint16 {
	return self.accessFlags
}
func (self *ClassFile) Fields() []*MemberInfo {
	return self.fields
}
func (self *ClassFile) Methods() []*MemberInfo {
	return self.methods
}


//以下3个 从常量池get相应的东西
func (self *ClassFile) ClassName() string {
	return self.constantPool.getClassName(self.thisClass)
}
func (self *ClassFile) SuperClassName() string {
	if self.superClass > 0 {
		return self.constantPool.getClassName(self.superClass)
	}
	return ""
}
func (self *ClassFile) InterfaceNames() []string {
	interfaceNames := make([]string, len(self.interfaces))
	for i, cpIndex := range self.interfaces {
		interfaceNames[i] = self.constantPool.getClassName(cpIndex)
	}
	return interfaceNames
}

func (self *ClassFile) SourceFileAttribute() *SourceFileAttribute {
	for _, attrInfo := range self.attributes {
		switch attrInfo.(type) {
		case *SourceFileAttribute:
			return attrInfo.(*SourceFileAttribute)
		}
	}
	return nil
}