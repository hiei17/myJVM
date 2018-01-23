package heap
//Class结构体，数组特有的方法

//类名[开头的就是数组类
//数组类不是加载的(米有.class文件 是虚拟机运行时产生的
func (self *Class) IsArray() bool {
	return self.name[0] == '['
}

//返回数组类的元素类型
func (self *Class) ComponentClass() *Class {
	//根据数组类名推测出数组元素类名，
	componentClassName := getComponentClassName(self.name)
	//类加载器加载元素类
	return self.loader.LoadClass(componentClassName)
}


//返回指定容量的数组对象
//通过类名能知道数组什么类型组成
func (self *Class) NewArray(count uint) *Object {
	if !self.IsArray() {
		panic("Not array class: " + self.name)
	}
	switch self.Name() {
		case "[Z":
			return &Object{self, make([]int8, count),nil}
		case "[B":
			return &Object{self, make([]int8, count),nil}
		case "[C":
			return &Object{self, make([]uint16, count),nil}
		case "[S":
			return &Object{self, make([]int16, count),nil}
		case "[I":
			return &Object{self, make([]int32, count),nil}
		case "[J":
			return &Object{self, make([]int64, count),nil}
		case "[F":
			return &Object{self, make([]float32, count),nil}
		case "[D":
			return &Object{self, make([]float64, count),nil}
		default:
			return &Object{self, make([]*Object, count),nil}
	}
}

//返回指定类的对应数组类
func (self *Class) ArrayClass() *Class {

	//组成的类 类名转变成类型描述符 前面加[ 就是数组类名
	arrayClassName := getArrayClassName(self.name)
	//加载数组类
	return self.loader.LoadClass(arrayClassName)
}