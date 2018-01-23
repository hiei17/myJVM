package heap

var primitiveTypes = map[string]string{
	"void":    "V",
	"boolean": "Z",
	"byte":    "B",
	"short":   "S",
	"int":     "I",
	"long":    "J",
	"char":    "C",
	"float":   "F",
	"double":  "D",
}

// [XXX -> [[XXX
// int -> [I
// XXX -> [LXXX;
func getArrayClassName(className string) string {
	return "[" + toDescriptor(className)
}

//根据数组类名推测出数组元素类名
// [[XXX -> [XXX
// [LXXX; -> XXX
// [I -> int
func getComponentClassName(className string) string {

	if className[0] == '[' {//数组类名以方括号开头
		//去掉[
		componentTypeDescriptor := className[1:]
		//然后把 类型描述符 转成 类名
		return toClassName(componentTypeDescriptor)
	}
	panic("Not array: " + className)
}

//类名转描述符
// [XXX => [XXX
// int  => I
// XXX  => LXXX;
func toDescriptor(className string) string {

	//数组类 类名就描述符
	if className[0] == '[' {
		// array
		return className
	}
	//基本类型名，返回对应的类型描述符
	if d, ok := primitiveTypes[className]; ok {
		// primitive
		return d
	}
	// object 普通类
	return "L" + className + ";"
}

//描述符转类名
// [XXX  => [XXX
// LXXX; => XXX
// I     => int
func toClassName(descriptor string) string {

	//肯定是数组，描述符即是类名。
	if descriptor[0] == '[' {
		// array
		return descriptor
	}

	//肯定的普通类  去掉开头的L和末尾的分号即是类名
	if descriptor[0] == 'L' {
		// object
		return descriptor[1 : len(descriptor)-1]
	}
	//是不是基本类型的描述符
	for className, d := range primitiveTypes {
		if d == descriptor {
			// primitive
			return className
		}
	}
	panic("Invalid descriptor: " + descriptor)
}
