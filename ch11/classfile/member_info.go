package classfile
//字段表 和 方法表 的公共部分

//数据结构
/*
field_info
{
	u2 access_flags;
	u2 name_index;
	u2 descriptor_index;
	u2 attributes_count;
	attribute_info attributes[attributes_count];
}
*/
type MemberInfo struct {
	cp ConstantPool//保存常量池指针
	accessFlags uint16 //访问标志
	nameIndex uint16//成员名(常量池坐标
	descriptorIndex uint16//描述符(常量池坐标
	attributes []AttributeInfo//属性表
}

//读取字段表或方法表
func readMembers(reader *ClassReader, cp ConstantPool) []*MemberInfo {
	memberCount := reader.readUint16()//前2个字节是个数
	members := make([]*MemberInfo, memberCount)//新建本结构体组成的数组
	for i := range members {
		///一个个成员读了放入数组
		members[i] = readMember(reader, cp)
	}
	return members
}

//返回本结构体
func readMember(reader *ClassReader, cp ConstantPool) *MemberInfo {

	return &MemberInfo{
			cp: cp,
			accessFlags: reader.readUint16(),//2字节访问标志
			nameIndex: reader.readUint16(),//2字节名字在常量池的索引
			descriptorIndex: reader.readUint16(),//2字节 描述符在常量池的索引
			attributes: readAttributes(reader, cp), // 成员属性 这里解析函数是通用的那个
	}

}

/******************************get方法 都不是给坐标 真的给值***********************************/
func (self *MemberInfo) AccessFlags() uint16 {
	return self.accessFlags;
}
func (self *MemberInfo) Name() string {
	return self.cp.getUtf8(self.nameIndex)
}
func (self *MemberInfo) Descriptor() string {
	return self.cp.getUtf8(self.descriptorIndex)
}
func (self *MemberInfo) CodeAttribute() *CodeAttribute {
	for _, attrInfo := range self.attributes {
		switch attrInfo.(type) {
		case *CodeAttribute:
			return attrInfo.(*CodeAttribute)
		}
	}
	return nil
}
func (self *MemberInfo) ConstantValueAttribute() *ConstantValueAttribute {
	for _, attrInfo := range self.attributes {
		switch attrInfo.(type) {
		case *ConstantValueAttribute:
			return attrInfo.(*ConstantValueAttribute)
		}
	}
	return nil
}
/***********************************************get方法***************************************************************/