package heap

import (
	"fmt"
	"jvmgo/ch11/classfile"
)

type Constant interface{}

//运行时常量池
type ConstantPool struct {
	class  *Class//所属类
	consts []Constant//常量数组
}

func newConstantPool(class *Class, cfCp classfile.ConstantPool) *ConstantPool {

	//new一个运行时常量池
	cpCount := len(cfCp)
	consts := make([]Constant, cpCount)
	rtCp := &ConstantPool{class, consts}

	//填consts
	for i := 1; i < cpCount; i++ {//遍历原classfile的常量池
		//常量数组
		cpInfo := cfCp[i]

		switch cpInfo.(type) {//常量类型

			//常量直接拿过来就好了
			case *classfile.ConstantIntegerInfo:
				//转型
				intInfo := cpInfo.(*classfile.ConstantIntegerInfo)
				consts[i] = intInfo.Value()//拿出值
			case *classfile.ConstantFloatInfo:
				floatInfo := cpInfo.(*classfile.ConstantFloatInfo)
				consts[i] = floatInfo.Value()
			case *classfile.ConstantLongInfo:
				longInfo := cpInfo.(*classfile.ConstantLongInfo)
				consts[i] = longInfo.Value()
				i++//占2个位
			case *classfile.ConstantDoubleInfo:
				doubleInfo := cpInfo.(*classfile.ConstantDoubleInfo)
				consts[i] = doubleInfo.Value()
				i++
			case *classfile.ConstantStringInfo:
				stringInfo := cpInfo.(*classfile.ConstantStringInfo)
				consts[i] = stringInfo.String()

			//引用那过来要加工下 不再放坐标了 直接实体放进去
			case *classfile.ConstantClassInfo://类信息//里面有常量池索引 指向类名
				classInfo := cpInfo.(*classfile.ConstantClassInfo)//强制向下转型
				//所属常量池  类名字符串
				consts[i] = newClassRef(rtCp, classInfo)

			case *classfile.ConstantFieldrefInfo://字段信息
				fieldrefInfo := cpInfo.(*classfile.ConstantFieldrefInfo)
				//所属常量池  	类名 自己的名字 描述符
				consts[i] = newFieldRef(rtCp, fieldrefInfo)

			case *classfile.ConstantMethodrefInfo://方法信息
				methodrefInfo := cpInfo.(*classfile.ConstantMethodrefInfo)
				//所属常量池  	类名 自己的名字 描述符
				consts[i] = newMethodRef(rtCp, methodrefInfo)

			case *classfile.ConstantInterfaceMethodrefInfo://接口信息
				methodrefInfo := cpInfo.(*classfile.ConstantInterfaceMethodrefInfo)
				//所属常量池  	类名 自己的名字 描述符
				consts[i] = newInterfaceMethodRef(rtCp, methodrefInfo)
			default:
			// todo
		}
	}

	return rtCp
}

func (self *ConstantPool) GetConstant(index uint) Constant {
	if c := self.consts[index]; c != nil {
		return c
	}
	panic(fmt.Sprintf("No constants at index %d", index))
}
